package weChatPay

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"yiarce/dragonnews/curl"
	"yiarce/dragonnews/general"
)

const (
	gateway     = "https://api.mch.weixin.qq.com"
	post        = "POST"
	get         = "GET"
	getCertList = "/v3/certificates"
)

const (
	a_256_GCM = "AEAD_AES_256_GCM"
)

type PayC struct {
	appId        string
	mchId        string
	appSecretKey string
	publicKey    *rsa.PublicKey
	privateKey   *rsa.PrivateKey
	serialNo     string
	notifyUrl    string
}

type CertList struct {
	Data []CertListData `json:"data"`
}

type CertListData struct {
	SerialNo           string            `json:"serial_no"`
	EffectiveTime      string            `json:"effective_time"`
	EncryptCertificate CertEnCertificate `json:"encrypt_certificate"`
}

type CertEnCertificate struct {
	Algorithm      string `json:"algorithm"`
	Nonce          string `json:"nonce"`
	AssociatedData string `json:"associated_data"`
	Ciphertext     string `json:"ciphertext"`
}

type Reply struct {
	Status     string
	StatusCode int
	RequestId  string
	Code       string
	Msg        string
	Body       *[]byte
	Request    *curl.Reply
}

//初始化微信支付对象
//  apiClientKey为微信证书所创建的密钥文件
//  传入对象可以使用ioutil.ReadFile("xxxx.pem")获得
//  apiClientCert为微信证书p12导出的pem格式的证书
//  传入对象可以使用ioutil.ReadFile("xxxx.pem")获得
//  notifyUrl将作为默认回调地址,提交请求时可临时更改
func Init(appId string, mchId string, secretKey string, apiClientKey string, apiClientCert string, notifyUrl string) (*PayC, error) {
	keyFile, err := ioutil.ReadFile(apiClientKey)
	if err != nil {
		return nil, err
	}
	rsaKey, err := decodeKey(keyFile)
	if err != nil {
		return nil, err
	}
	certFile, err := ioutil.ReadFile(apiClientCert)
	if err != nil {
		return nil, err
	}
	serialNo, publicKey, err := decodeCert(certFile)
	if err != nil {
		return nil, err
	}
	return &PayC{
		appId:        appId,
		mchId:        mchId,
		appSecretKey: secretKey,
		publicKey:    publicKey,
		privateKey:   rsaKey,
		serialNo:     serialNo,
		notifyUrl:    notifyUrl,
	}, nil
}

func (p *PayC) App() *app {
	return &app{p: p}
}

//生成验证请求头内容,可用于第三方工具调用验证是否有误
func Authorization(sign string, mchId string, timestamp string, nonce string, serialNo string) string {
	return `WECHATPAY2-SHA256-RSA2048 mchid="` + mchId + `",nonce_str="` + nonce + `" ,signature="` + sign + `",timestamp="` + timestamp + `",serial_no="` + serialNo + `"`
}

//内部生成验证请求头
func (p *PayC) header(sign string, timestamp string, nonce string) map[string]string {
	return map[string]string{
		"Authorization": `WECHATPAY2-SHA256-RSA2048 mchid="` + p.mchId + `",nonce_str="` + nonce + `" ,signature="` + sign + `",timestamp="` + timestamp + `",serial_no="` + p.serialNo + `"`,
		"User-Agent":    `DragonNews/1.0 (` + runtime.GOOS + `; ` + runtime.GOARCH + `) ` + ` ` + runtime.Version(),
		"Accept":        "application/json",
	}
}

//生成签名
func (p *PayC) createSign(method string, url string, timestamp string, nonce string, text string) string {
	data := method + "\n" + url + "\n" + timestamp + "\n" + nonce + "\n" + text + "\n"
	hashd := sha256.Sum256([]byte(data))
	sign, _ := p.privateKey.Sign(rand.Reader, hashd[:], crypto.SHA256)
	return base64.StdEncoding.EncodeToString(sign)
}

//输出签名,此生成的签名可用于官方验签工具验证是否正确
func CreateSign(text string, apiClientKey []byte) (string, error) {
	pri, err := decodeKey(apiClientKey)
	if err != nil {
		return "", err
	}
	hashd := sha256.Sum256([]byte(text))
	sign, _ := pri.Sign(rand.Reader, hashd[:], crypto.SHA256)
	return string(sign), nil
}

//下载证书
func DownloadCert() {

}

//获取证书列表
func (p *PayC) GetCertList(nonce string) (Reply, map[int]interface{}, error) {
	r, err := p.send(get, getCertList, general.Date().Timestamp("s"), nonce, nil)
	if err != nil {
		return r, map[int]interface{}{}, err
	}
	cert := CertList{}
	json.Unmarshal(r.Request.Body, &cert)
	certData := map[int]interface{}{}
	for key, val := range cert.Data {
		pla, err := decodeData(val.EncryptCertificate.Ciphertext, val.EncryptCertificate.Nonce, "ycwlf101ycwlf101ycwlf101ycwlf101", val.EncryptCertificate.AssociatedData, val.EncryptCertificate.Algorithm)
		if err != nil {
			return r, map[int]interface{}{}, errors.New("数据解密失败,可自行解密!")
		}
		certData[key] = pla
	}
	return r, certData, err
}

//提交请求
func (p *PayC) send(method string, url string, timestamp string, nonce string, data []byte) (Reply, error) {
	var c curl.Reply
	var err error
	if method == post {
		c, err = curl.Post(gateway+url, curl.Json, p.header(p.createSign(post, url, timestamp, nonce, string(data)), timestamp, nonce), string(data))
	} else {
		c, err = curl.Get(gateway+url, curl.Html, p.header(p.createSign(get, url, timestamp, nonce, ""), timestamp, nonce))
	}
	r := Reply{}
	r.Request = &c
	if err != nil {
		return r, err
	}
	r.RequestId = c.Header["Request-Id"][0]
	switch c.StatusCode {
	case 200, 204:
		r.StatusCode = c.StatusCode
		r.Status = c.Status
		r.Body = &c.Body
		break
	case 202:
		return p.send(method, url, timestamp, nonce, data)
	default:
		r.Status = c.Status
		r.StatusCode = c.StatusCode
		d := map[string]string{}
		json.Unmarshal(c.Body, &d)
		r.Code = d["code"]
		r.Msg = d["message"]
		break
	}
	//if c.Header["Wechatpay-Serial"][0] != p.serialNo {
	//	return r, errors.New("证书不匹配,请更换证书!")
	//}
	//if p.verifiedSign(c.Header["Wechatpay-Signature"][0],c.Header["Wechatpay-Timestamp"][0],c.Header["Wechatpay-Nonce"][0],c.Body) != nil {
	//	return r, errors.New("非微信官方的应答,签名验证失败!")
	//}
	return r, nil
}

//验证签名
func (p *PayC) verifiedSign(sign string, timestamp string, nonce string, text []byte) error {
	data := timestamp + "\n" + nonce + "\n" + string(text) + "\n"
	hashd := sha256.Sum256([]byte(data))
	err := rsa.VerifyPKCS1v15(p.publicKey, crypto.SHA256, hashd[:], []byte(sign))
	return err
}

//获得文件密钥
func decodeKey(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pri, flag := key.(*rsa.PrivateKey)
	if !flag {
		return nil, errors.New("非对应密钥,请检查密钥文件")
	}
	return pri, nil
}

//解密证书
func decodeCert(bytes []byte) (string, *rsa.PublicKey, error) {
	block, _ := pem.Decode(bytes)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", nil, err
	}
	rsaPub, flag := cert.PublicKey.(*rsa.PublicKey)
	if !flag {
		return "", nil, err
	}
	return strings.ToUpper(hex.EncodeToString(cert.SerialNumber.Bytes())), rsaPub, err
}

//解密数据包
func decodeData(cip string, nonce string, key string, associated string, enType string) (string, error) {
	switch enType {
	case a_256_GCM:
		p, err := aes.NewCipher([]byte(key))
		if err != nil {
			return "", err
		}
		aead, aeadErr := cipher.NewGCM(p)
		if aeadErr != nil {
			return "", aeadErr
		}
		decip, _ := base64.StdEncoding.DecodeString(cip)
		pla, dErr := aead.Open(nil, []byte(nonce), decip, []byte(associated))
		if dErr != nil {
			fmt.Println(dErr)
			return "", dErr
		}
		return string(pla), nil
	default:
		return "", errors.New("未识别的解密方式: " + enType)
	}
}
