package cmdAnime

import (
	"fmt"
	"yiarce/core/date"
)

type LoadS struct {
	Str  string
	P    string
	ls   string
	flag bool
	time int
}

func LoadAnime(l *LoadS, flag *bool) bool {
	if l.P == "" {
		l.P = "."
		l.ls = l.P
	}
	if l.flag == false {
		fmt.Println(l.Str + l.ls)
		l.flag = true
	}
	if !(*flag) {
		Clear()
		fmt.Println(l.Str + l.ls)
		if date.Date().Unix()-l.time > 1 {
			l.time = date.Date().Unix()
			if len(l.ls) == (len(l.P) * 3) {
				l.ls = l.P
			} else {
				l.ls += l.P
			}
		}
		return true
	}
	fmt.Println(`我退出了`)
	return false
}

func Clear() {
	fmt.Printf("\033[1A\033[K")
}
