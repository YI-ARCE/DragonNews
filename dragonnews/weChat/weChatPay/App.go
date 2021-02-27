package weChatPay

import (
	"encoding/json"
)

const (
	appUniformOrder = "/v3/pay/transactions/app"
)

type app struct {
	p *PayC
	*generals
}

//统一下单提交数据
type AppUniformOrderData struct {
	AppId       string  `json:"appid"`
	MchId       string  `json:"mchid"`
	Description string  `json:"description"`
	OutTradeNo  string  `json:"out_trade_no"`
	TimeExpire  string  `json:"time_expire,omitempty"`
	Attach      string  `json:"attach"`
	NotifyUrl   string  `json:"notify_url"`
	GoodsTag    string  `json:"goods_tag"`
	Amount      Amount  `json:"amount"`
	Detail      *Detail `json:"detail,omitempty"`
	SceneInfo   *Secene `json:"scene_info,omitempty"`
}

//统一下单
func (g *app) UniformOrder(data AppUniformOrderData, timestamp string, nonce string) (Reply, error) {
	data.AppId = g.p.appId
	data.MchId = g.p.mchId
	if data.NotifyUrl == "" {
		data.NotifyUrl = g.p.notifyUrl
	}
	str, _ := json.Marshal(data)
	return g.p.send(post, appUniformOrder, timestamp, nonce, str)
}

//只提交必交字段,若无其它硬性要求,可使用此方法构建基础内容
func (u *AppUniformOrderData) Default(total int, tradeName string, tradeNo string, notifyUrl ...string) {
	u.Amount.Total = total
	u.Description = tradeName
	u.OutTradeNo = tradeNo
	if len(notifyUrl) > 0 {
		u.NotifyUrl = notifyUrl[0]
	}
	u.Detail = nil
	u.SceneInfo = nil
}
