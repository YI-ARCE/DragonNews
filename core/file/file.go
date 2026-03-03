package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type File struct {
	self     *os.File
	Index    int64
	readBuff []byte
}

type Instant struct {
	path     string
	filename string
	mode     os.FileMode
}

func GetCustom(filename string, mode int) (*File, error) {
	file, err := os.OpenFile(filename, mode, 0777)
	if err != nil {
		return nil, err
	}
	return &File{file, 0, nil}, nil
}

func Get(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &File{f, 0, nil}, nil
}

// Set 文件保存
//
//	path 存储的路径
//	filename 保存的文件名
//	flag os包下的O_开头的常量
//	mode unix下的读写权限,格式0777
func Set(path string, filename string, content []byte, flag int, mode os.FileMode) error {
	// 验证路径和文件名参数
	if path == "" || filename == "" {
		return fmt.Errorf("path and filename cannot be empty")
	}

	// 创建目录
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	// 创建或打开文件
	f, err := os.OpenFile(path+`/`+filename, flag, mode)
	if err != nil {
		return err
	}

	// 确保文件关闭
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// 记录错误，但不影响主函数的返回值
			fmt.Println("Error closing file:", closeErr.Error())
		}
	}()

	// 写入内容
	if content != nil {
		_, err = f.Write(content)
		if err != nil {
			return err
		}
		// 强制刷新缓冲区
		if syncErr := f.Sync(); syncErr != nil {
			// 记录错误，但不影响主函数的返回值
			fmt.Println("Error syncing file:", syncErr.Error())
		}
	}

	// 清空content，帮助GC
	content = nil
	return nil
}

func (f *File) String() string {
	r, err := io.ReadAll(f.self)
	if err != nil {
		return ""
	}
	return string(r)
}

func (f *File) Byte() []byte {
	r, err := io.ReadAll(f.self)
	if err != nil {
		return []byte{}
	}
	return r
}

func (f *File) GetReader() *os.File {
	return f.self
}

func (f *File) Read(length int64) ([]byte, error) {
	// 验证文件对象是否为nil
	if f == nil || f.self == nil {
		return nil, fmt.Errorf("file object is nil")
	}

	// 验证长度参数
	if length <= 0 {
		return []byte{}, nil
	}

	// 分配缓冲区
	b := make([]byte, length)

	// 读取数据
	n, err := f.self.ReadAt(b, f.Index)
	if err != nil {
		if err == io.EOF {
			// 到达文件末尾，返回已读取的数据
			return b[:n], nil
		}
		// 其他错误
		return nil, err
	}

	// 更新读取位置
	f.Index += int64(n)

	return b[:n], nil
}

func (f *File) ReadLine(function func(raw string)) {
	// 验证文件对象和回调函数是否为nil
	if f == nil || f.self == nil || function == nil {
		return
	}

	// 初始化读取缓冲区
	if f.readBuff == nil {
		f.readBuff = []byte{}
	}

	// 换行符
	scan := []byte("\n")

	for {
		// 读取数据
		read, err := f.Read(4096)
		if err != nil {
			// 处理错误
			fmt.Println("ReadLine error:", err.Error())
			break
		}

		// 检查是否到达文件末尾
		if len(read) == 0 {
			// 处理最后一行没有换行符的情况
			if len(f.readBuff) > 0 {
				function(string(f.readBuff))
				f.readBuff = []byte{}
			}
			break
		}

		// 追加读取的数据到缓冲区
		f.readBuff = append(f.readBuff, read...)

		// 处理缓冲区中的行
		for {
			index := bytes.Index(f.readBuff, scan)
			if index == -1 {
				break
			}

			// 提取一行数据
			str := string(f.readBuff[:index])
			// 更新缓冲区
			f.readBuff = f.readBuff[index+1:]
			// 调用回调函数
			function(str)
		}
	}
}

func New(path string, filename string, mode os.FileMode) (*File, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path+`/`+filename, os.O_TRUNC|os.O_CREATE, mode)
	if err != nil {
		return nil, err
	}
	return &File{
		self:  f,
		Index: 0,
	}, nil
}

func (f *File) Append(content string) error {
	// 验证文件对象是否为nil
	if f == nil || f.self == nil {
		return fmt.Errorf("file object is nil")
	}

	// 写入内容
	_, err := f.self.WriteString(content)
	return err
}

func (f *File) Close() error {
	// 验证文件对象是否为nil
	if f == nil || f.self == nil {
		return fmt.Errorf("file object is nil")
	}

	// 关闭文件
	return f.self.Close()
}

func NewInstant(path string, filename string, str string, mode os.FileMode) *Instant {
	// 验证路径和文件名参数
	if path == "" || filename == "" {
		return nil
	}

	// 创建文件
	err := Set(path, filename, []byte(str), os.O_TRUNC|os.O_CREATE, mode)
	if err != nil {
		// 记录错误，但仍然返回Instant对象
		fmt.Println("Error creating file:", err.Error())
	}

	return &Instant{
		path:     path,
		filename: filename,
		mode:     mode,
	}
}

func (i *Instant) Append(content string) error {
	// 验证Instant对象是否为nil
	if i == nil || i.path == "" || i.filename == "" {
		return fmt.Errorf("instant object is nil or path/filename is empty")
	}

	// 打开文件进行追加
	f, err := os.OpenFile(i.path+`/`+i.filename, os.O_APPEND|os.O_WRONLY, i.mode)
	if err != nil {
		return err
	}

	// 确保文件关闭
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// 记录错误，但不影响主函数的返回值
			fmt.Println("Error closing file:", closeErr.Error())
		}
	}()

	// 写入内容
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	// 清空content，帮助GC
	content = ``
	return nil
}
