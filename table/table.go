package table

import "strings"

//type defaults struct {
//	Name  string
//	Alias string
//}

func String(as *string, s *[]string) string {
	str := strings.Builder{}
	l := len(*s)
	for i := 0; i <= l; i++ {
		str.WriteString("`")
		str.WriteString(*as)
		str.WriteString("`.")
		str.WriteString((*s)[i])
		if i != l {
			str.WriteString(`,`)
		}
	}
	return str.String()
}
