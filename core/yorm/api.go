package yorm

func Raw(str string) string {
	return `[raw]__dn:` + str
}
