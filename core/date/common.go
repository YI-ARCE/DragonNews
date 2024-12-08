package date

import (
	"strconv"
	"time"
)

type date struct {
	d time.Time
}

const DATE = "2006-01-02"

const DATETIME = "2006-01-02 15:04:05"

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

func Date() date {
	d := time.Now()
	return date{d: d}
}

func DateTime() string {
	d := time.Now()
	d.Format("Y")
	return d.Format("2006-01-02 15:04:05")
}

func ParseDate(dates string) date {
	parse, err := time.Parse("2006-01-02 15:04:05", dates)
	if err != nil {
		panic(err.Error())
	}
	return date{d: parse}
}

func CheckDate(dates string, format string) error {
	_, err := time.Parse(format, dates)
	if err != nil {
		return err
	}
	return nil
}

// Timestamp types为输出格式,填入s为秒级时间戳,ms毫秒,ns纳秒
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

func (D date) Source() time.Time {
	return D.d
}

// Year 输出完整的年份 如2021
func (D date) Year() string {
	return D.d.Format("2006")
}

// Month 输出完整的月份 如01
func (D date) Month() string {
	return D.d.Format("01")
}

// Day 输出补位的天数 如01
func (D date) Day() string {
	return D.d.Format("02")
}

// Hour 输出补位的24小时制 如09
func (D date) Hour() string {
	return D.d.Format("15")
}

// Minutes 输出补位的分钟 如00
func (D date) Minutes() string {
	return D.d.Format("04")
}

// Seconds 输出补位的秒 如00
func (D date) Seconds() string {
	return D.d.Format("05")
}

// YM 输出年月组合,根据传入的格式符输出 如传入‘-’,则输出"2021-05"
func (D date) YM(d string) string {
	return D.d.Format("2006" + d + "01")
}

// YMD 输出年月日组合,根据传入的格式符输出 如传入‘-’,则输出"2021-05-01"
func (D date) YMD(d string) string {
	return D.d.Format("2006" + d + "01" + d + "02")
}

// HIS 输出时分秒组合,根据传入的格式符输出 如传入‘:’,则输出"09:05:01"
func (D date) HIS(d string) string {
	return D.d.Format("15" + d + "04" + d + "05")
}

func (D date) Unix() int {
	return int(D.d.Unix())
}

// Custom Y M D H I S 代表补位的年月日时分秒,Y则为完整年份,例:2021
//
//	y m d h i s 代表小于 10 的不补 0,y则为当前年的两位数,例:2021->21
//	可自行调整输出格式,例:Y-M-D H:I:S 则会输出"2021-01-01 09:01:01"此类格式
func (D date) Custom(s ...string) string {
	format := ""
	if len(s) > 0 {
		format = s[0]
	} else {
		format = "Y-M-D H:I:S"
	}
	dt := ""
	Len := len(format) - 1
	for i := 0; i <= Len; i++ {
		index := format[i : i+1]
		if index == "H" || index == "h" {
			predict := format[i : i+2]
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
			dt += format[i : i+1]
		}
	}
	return dt
}

func Time(t int64) date {
	tm := time.Unix(t, 0)
	return date{d: tm}
}

func TimeMill(t int64) date {
	tm := time.UnixMilli(t)
	return date{d: tm}
}
