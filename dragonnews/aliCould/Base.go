package aliCould

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type Client struct {
	C *dysmsapi.Client
}

func Init(accessKey string, accessKeySecret string) (Client, error) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKey, accessKeySecret)
	if err != nil {
		return Client{}, err
	}
	return Client{client}, nil
}

func (c Client) SendSms(phone string, body string, singName string, code string) (*dysmsapi.SendSmsResponse, error) {
	r := dysmsapi.CreateSendSmsRequest()
	r.PhoneNumbers = phone
	r.SignName = singName
	r.TemplateCode = code
	r.TemplateParam = body
	rs, err := c.C.SendSms(r)
	if err != nil {
		return &dysmsapi.SendSmsResponse{}, err
	}
	return rs, nil
}
