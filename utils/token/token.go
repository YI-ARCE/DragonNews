package token

import (
	"math/rand"
	"strconv"
	"strings"
	"yiarce/core/date"
	"yiarce/core/encrypt"
	"yiarce/core/frame"
)

const key = `yiarceloveSakura`

const _tag = `token`

// Context 令牌结构
// 包含用户ID、手机号、登录时间和过期时间
type Context struct {
	ID          int    //ID
	Phone       string // 手机号
	LoginTime   int    // 登录时间
	ExpiredTime int    // 过期时间
}

// CreateToken 创建用户令牌
//
// 参数：
//   - id: 用户ID
//   - phone: 用户手机号
//
// 返回值：
//   - string: 生成的token字符串
//
// 说明：
//   - 使用AES加密生成token
//   - token包含用户ID、随机数、手机号、登录时间和过期时间
//   - 过期时间为登录时间后15天（1296000秒）
func CreateToken(id int, phone string) string {
	// 检查参数
	if id <= 0 {
		return ""
	}

	if phone == "" {
		return ""
	}

	times := date.New().Unix()
	// 创建AES加密对象
	aesObj, err := encrypt.Aes(key)
	if err != nil {
		return ""
	}

	// 加密
	encrypted, err := aesObj.Encrypt(strconv.Itoa(id) + `,` + strconv.Itoa(rand.Intn(9999999)) + `,` + phone + `,` + strconv.Itoa(times) + `,` + strconv.Itoa(times+1296000))
	if err != nil {
		return ""
	}

	return encrypted.ToBase64()
}

// DecryptToken 解密用户令牌
//
// 参数：
//   - cip: 加密的token字符串
//
// 返回值：
//   - Token: 解密后的token结构
//   - error: 错误信息
//
// 说明：
//   - 使用AES解密token
//   - 验证token格式和过期时间
//   - 返回包含用户信息的Token结构
func DecryptToken(cip string) (Context, error) {
	if len(cip) < 1 {
		return Context{}, frame.NewError(_tag, `未登录用户`)
	}

	// 创建AES加密对象
	aesObj, err := encrypt.Aes(key)
	if err != nil {
		return Context{}, frame.NewError(frame.SelfError, "Failed to create AES cipher: "+err.Error())
	}

	// 解密
	decryptedObj, err := aesObj.Decrypt(cip)
	if err != nil {
		return Context{}, frame.NewError(_tag, "Failed to decrypt token: "+err.Error())
	}

	decrypted := decryptedObj.ToString()
	arr := strings.Split(decrypted, `,`)
	if len(arr) < 5 {
		return Context{}, frame.NewError(_tag, `无效的token格式`)
	}

	id, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return Context{}, frame.NewError(_tag, `无效的token格式`)
	}

	loginTime, err := strconv.ParseInt(arr[3], 10, 64)
	if err != nil {
		return Context{}, frame.NewError(_tag, `无效的token格式`)
	}

	expireTime, err := strconv.ParseInt(arr[4], 10, 64)
	if err != nil {
		return Context{}, frame.NewError(_tag, `无效的token格式`)
	}

	// 检查token是否过期
	if int(expireTime) < date.New().Unix() {
		return Context{}, frame.NewError(_tag, `token已过期`)
	}

	return Context{
		ID:          int(id),
		Phone:       arr[2],
		LoginTime:   int(loginTime),
		ExpiredTime: int(expireTime),
	}, nil
}
