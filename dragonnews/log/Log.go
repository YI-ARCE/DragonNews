package log

import (
	"errors"
	"fmt"
	"os"
	"yiarce/dragonnews/general"
	. "yiarce/dragonnews/log/driver"
)

type Log struct {
	Host   string
	Method string
	Uri    string
	IP     string
	flag   bool
	msg    string
	end    string
}

func Init(Host string, Method string, Uri string, IP string) Log {
	l := Log{}
	l.Host = Host
	l.Method = Method
	l.Uri = Uri
	l.IP = IP
	return l
}

func (l *Log) Insert(s interface{}) {
	str, err := Serialize(s)
	if err != nil {
		return
	}
	if len(l.msg) < 1 {
		head := "----[" + general.Date().Custom("Y-M-D H:I:S") + "]----[" + l.IP + "]----[" + l.Method + "]----[" + l.Host + "]----[" + l.Uri + "]----"
		l.msg = head + "\n◎ " + str + "\n"
		Len := len(head)
		for i := 0; i < Len; i++ {
			l.end += "-"
		}
		l.end += "\n"
	} else {
		l.msg += "◎ " + str + "\n"
	}
}

func (l *Log) Judge() bool {
	if l.msg != "" {
		return true
	}
	return false
}

func Out(l Log) error {
	var path, _ = os.Getwd()
	d := general.Date()
	path += "\\log\\" + d.Custom("YM") + "\\"
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println(err)
		return errors.New("创建日志文件失败!")
	}
	filename := d.Custom("D") + ".txt"
	file, err := os.OpenFile(path+filename, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		return errors.New("创建日志文件失败!")
	}
	_, err = file.WriteString(l.msg + l.end)
	if err != nil {
		return errors.New("日志文件数据写入失败!")
	}
	file.Close()
	return nil
}

func String() {

}
