package weChat

import (
	"encoding/json"
	"errors"
	"yiarce/dragonnews/curl"
)

const (
	gateway          = "https://api.weixin.qq.com"
	oauth2AccessT    = "/sns/oauth2/access_token"
	oauth2RefAccessT = "/sns/oauth2/refresh_token"
)

type wx struct {
	appid     string
	secret    string
	key       string
	mchId     string
	notifyUrl string
}

type Oauth2 struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
	Scope        string `json:"scope"`
	ErrMsg       string `json:"errmsg"`
	ErrCode      int    `json:"errcode"`
}

//初始化
func Init(appid string, appSecret string, apiKey string, mchId string, notifyUrl string) *wx {
	return &wx{appid, appSecret, apiKey, mchId, notifyUrl}
}

//用户授权,当参数二不为空时则作为刷新token使用
//  code为授权码,当参数二不为nil时code参数应传入refresh_token
func (w *wx) Oauth(code string, flag ...int) (Oauth2, error) {
	query := ""
	if len(flag) > 0 {
		query = oauth2RefAccessT + "?appid=" + w.appid + "&grant_type=refresh_token&refresh_token=" + code
	} else {
		query = oauth2AccessT + "?appid=" + w.appid + "&secret=" + w.secret + "&code=" + code + "&grant_type=authorization_code"
	}
	r, err := w.send(query, nil, nil)
	if err != nil {
		return Oauth2{}, err
	}
	data := Oauth2{}
	json.Unmarshal(r.Body, &data)
	if data.ErrCode != 0 {
		return data, errors.New("请求失败,请打印返回的数据查看!")
	}
	return data, nil
}

//发送
func (w *wx) send(url string, header map[string]string, data []byte) (curl.Reply, error) {
	if data != nil {
		return curl.Post(gateway+url, curl.Json, header, string(data))
	}
	return curl.Get(gateway+url, curl.Html, header)
}
