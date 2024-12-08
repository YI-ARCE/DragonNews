package dhttp

import (
	"math/rand"
	"strconv"
	"strings"
	"yiarce/core/date"
	"yiarce/core/encrypt"
	"yiarce/core/frame"
)

type Token struct {
	Id          int    //ID
	Phone       string // 手机号
	LoginTime   int    // 登录时间
	ExpiredTime int    // 过期时间
}

func CreateToken(id int, phone string) string {
	times := date.Date().Unix()
	return encrypt.Aes(`yiarceLoveSakura`).Encrypt(strconv.Itoa(id) + `,` + strconv.Itoa(rand.Intn(9999999)) + `,` + phone + `,` + strconv.Itoa(times) + `,` + strconv.Itoa(times+1296000)).ToBase64()
}

func DecryptToken(cip string) Token {
	if len(cip) < 1 {
		frame.Errors(frame.TokenError, `未登录用户`, nil)
	}
	defer frame.Prevent(frame.TokenError, `1003`)
	arr := strings.Split(encrypt.Aes(`yiarceLoveSakura`).Decrypt(cip).ToString(), `,`)
	id, _ := strconv.ParseInt(arr[0], 10, 64)
	loginTime, _ := strconv.ParseInt(arr[3], 10, 64)
	expireTime, _ := strconv.ParseInt(arr[4], 10, 64)
	return Token{
		int(id),
		arr[2],
		int(loginTime),
		int(expireTime),
	}
}
