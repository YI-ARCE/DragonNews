package weChatPay

import (
	"encoding/json"
)

const (
	h5UniformOrder = "/v3/pay/transactions/h5"
)

type h5 struct {
	p *PayC
	*generals
}

//统一下单提交数据
type H5UniformOrderData struct {
	AppId       string      `json:"appid"`
	MchId       string      `json:"mchid"`
	Description string      `json:"description"`
	OutTradeNo  string      `json:"out_trade_no"`
	TimeExpire  string      `json:"time_expire,omitempty"`
	Attach      string      `json:"attach"`
	NotifyUrl   string      `json:"notify_url"`
	GoodsTag    string      `json:"goods_tag"`
	Amount      Amount      `json:"amount"`
	Detail      *Detail     `json:"detail,omitempty"`
	SceneInfo   *H5Secene   `json:"scene_info,omitempty"`
	SettleInfo  *SettleInfo `json:"settle_info,omitempty"`
}

//场景信息
type H5Secene struct {
	PayerClientIp string    `json:"payer_client_ip"`
	DeviceId      string    `json:"device_id"`
	StoreInfo     StoreInfo `json:"store_info"`
	H5Info        H5Info    `json:"h5_info"`
}

type H5Info struct {
	Type        string `json:"type"`
	AppName     string `json:"app_name"`
	AppUrl      string `json:"app_url"`
	BundleId    string `json:"bundle_id"`
	PackageName string `json:"package_name"`
}

//统一下单
func (h *h5) UniformOrder(data H5UniformOrderData, timestamp string, nonce string) (Reply, error) {
	data.AppId = h.p.appId
	data.MchId = h.p.mchId
	if data.NotifyUrl == "" {
		data.NotifyUrl = h.p.notifyUrl
	}
	str, _ := json.Marshal(data)
	return h.p.send(post, h5UniformOrder, timestamp, nonce, str)
}

//只提交必交字段,若无其它硬性要求,可使用此方法构建基础内容
func (h *H5UniformOrderData) Default(total int, tradeName string, tradeNo string, notifyUrl ...string) {
	h.Amount.Total = total
	h.Description = tradeName
	h.OutTradeNo = tradeNo
	if len(notifyUrl) > 0 {
		h.NotifyUrl = notifyUrl[0]
	}
	h.Detail = nil
	h.SceneInfo = nil
	h.SettleInfo = nil
}
