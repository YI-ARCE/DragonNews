package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unsafe"
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
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path+`/`+filename, flag, mode)
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	content = nil
	return nil
}

func (f *File) String() string {
	r, _ := io.ReadAll(f.self)
	return *(*string)(unsafe.Pointer(&r))
}

func (f *File) Byte() []byte {
	r, _ := io.ReadAll(f.self)
	return r
}

func (f *File) GetReader() *os.File {
	return f.self
}

func (f *File) Read(length int64) ([]byte, error) {
	b := make([]byte, length)
	n, err := f.self.ReadAt(b, f.Index)
	if err != nil {
		if err == io.EOF {
			return b[:n], err
		} else {
			return nil, err
		}
	}
	f.Index += length
	return b, err
}

func (f *File) ReadLine(function func(raw string)) {
	if f.readBuff == nil {
		f.readBuff = []byte{}
	}
	scan := []byte("\n")
	for {
		read, err := f.Read(4096)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err.Error())
			}
		}
		f.readBuff = append(f.readBuff, read...)
		for {
			index := bytes.Index(f.readBuff, scan)
			if index == -1 {
				break
			}
			str := string(f.readBuff[:index])
			f.readBuff = f.readBuff[index+1:]
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
	_, err := f.self.WriteString(content)
	return err
}

func (f *File) Close() error {
	return f.self.Close()
}

func NewInstant(path string, filename string, str string, mode os.FileMode) *Instant {
	Set(path, filename, []byte(str), os.O_TRUNC|os.O_CREATE, mode)
	return &Instant{
		path:     path,
		filename: filename,
		mode:     mode,
	}
}

func (i *Instant) Append(content string) error {
	f, err := os.OpenFile(i.path+`/`+i.filename, os.O_APPEND, i.mode)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	f.Close()
	content = ``
	return nil
}
