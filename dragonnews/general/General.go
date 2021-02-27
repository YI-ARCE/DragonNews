package general

import (
	"strconv"
	"time"
)

//const (
//	Y  = "2006"
//	y  = "06"
//	M  = "01"
//	m  = "1"
//	D  = "02"
//	d  = "2"
//	H  = "15"
//	HH = ""
//	h  = "3"
//	hh = "03"
//	I  = "04"
//	i  = "4"
//	S  = "05"
//	s  = "5"
//)

type date struct {
	d time.Time
}

var arr = map[string]string{
	"Y":  "2006",
	"y":  "06",
	"M":  "01",
	"m":  "1",
	"D":  "02",
	"d":  "2",
	"H":  "15",
	"HH": "",
	"h":  "3",
	"hh": "03",
	"I":  "04",
	"i":  "4",
	"S":  "05",
	"s":  "5",
}

// Y代表完整年份,y代表两位数年份
//
// M代表月份不补0,m代表补0
//
// D代表天不补0,d代表补0
//
// HH代表12小时制不补0,hh代表12小时制补0
//
// H代表24小时制不补0,h代表24小时制补0
//
// I代表分不补0,i代表分补0
//
// S代表秒不补0,s代表补0 , 例子:
//  Y-m-d h:i:s 解释为2021-01-25 16:01:01
//  y-M-D H:I:S 解释为21-1-25 16:1:1
//  y-M-D HH:I:S 解释为21-1-25 4:1:1
//  y-M-D hh:I:S 解释为21-1-25 04:1:0
func Date() date {
	d := time.Now()
	return date{d: d}
}

func DateTime() string {
	d := time.Now()
	return d.Format("2006-01-02 15:04:05")
}

//types为输出格式,填入s为秒级时间戳,ms毫秒,ns纳秒
func (D date) Timestamp(types string) string {
	t := ""
	switch types {
	case "s":
		t = strconv.FormatInt(D.d.Unix(), 10)
		break
	case "ms":
		t = strconv.FormatInt(D.d.UnixNano()/1e6, 10)
		break
	case "ns":
		t = strconv.FormatInt(D.d.UnixNano(), 10)
		break
	default:
		break
	}
	return t
}

func (D date) Y() string {
	return strconv.Itoa(D.d.Year())
}

func (D date) Ym(d string) string {
	return strconv.Itoa(D.d.Year()) + d + strconv.Itoa(int(D.d.Month()))
}

func (D date) Ymd(d string) string {
	return strconv.Itoa(D.d.Year()) + d + strconv.Itoa(int(D.d.Month())) + d + strconv.Itoa(D.d.Day())
}

func (D date) D() string {
	return strconv.Itoa(D.d.Day())
}

func (D date) Custom(s string) string {
	dt := ""
	Len := len(s) - 1
	for i := 0; i <= Len; i++ {
		index := s[i : i+1]
		if index == "H" || index == "h" {
			predict := s[i : i+2]
			if predict == "HH" || predict == "hh" {
				index = predict
				i++
			}
		}
		if val, flag := arr[index]; flag {
			if index == "HH" {
				dt += strconv.Itoa(D.d.Hour())
			} else {
				dt += D.d.Format(val)
			}
		} else {
			dt += s[i : i+1]
		}
	}
	return dt
}
