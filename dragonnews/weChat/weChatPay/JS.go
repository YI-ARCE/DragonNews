package weChatPay

import (
	"encoding/json"
)

const (
	jsUniformOrder = "/v3/pay/transactions/jsapi"
)

type js struct {
	p *PayC
	*generals
}

//统一下单提交数据
type JsUniformOrderData struct {
	AppId       string      `json:"appid"`
	MchId       string      `json:"mchid"`
	Description string      `json:"description"`
	OutTradeNo  string      `json:"out_trade_no"`
	TimeExpire  string      `json:"time_expire,omitempty"`
	Attach      string      `json:"attach"`
	NotifyUrl   string      `json:"notify_url"`
	GoodsTag    string      `json:"goods_tag"`
	Amount      Amount      `json:"amount"`
	Paper       *JsPaper    `json:"paper"`
	Detail      *Detail     `json:"detail,omitempty"`
	SceneInfo   *Secene     `json:"scene_info,omitempty"`
	SettleInfo  *SettleInfo `json:"settle_info,omitempty"`
}

type JsPaper struct {
	OpenId string `json:"openid"`
}

//统一下单
func (j *js) UniformOrder(data JsUniformOrderData, timestamp string, nonce string) (Reply, error) {
	data.AppId = j.p.appId
	data.MchId = j.p.mchId
	if data.NotifyUrl == "" {
		data.NotifyUrl = j.p.notifyUrl
	}
	str, _ := json.Marshal(data)
	return j.p.send(post, jsUniformOrder, timestamp, nonce, str)
}

//只提交必交字段,若无其它硬性要求,可使用此方法构建基础内容
func (j *JsUniformOrderData) Default(total int, tradeName string, openId string, tradeNo string, notifyUrl ...string) {
	j.Amount.Total = total
	j.Paper.OpenId = openId
	j.Description = tradeName
	j.OutTradeNo = tradeNo
	if len(notifyUrl) > 0 {
		j.NotifyUrl = notifyUrl[0]
	}
	j.Detail = nil
	j.SceneInfo = nil
	j.SettleInfo = nil
}
