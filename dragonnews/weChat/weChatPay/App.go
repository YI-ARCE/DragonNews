package weChatPay

import (
	"encoding/json"
	"yiarce/dragonnews/general"
)

const (
	uniformOrder  = "/v3/pay/transactions/app"
	queryTradeNo  = "/v3/pay/transactions/out-trade-no/"
	applyRefund   = "/v3/refund/domestic/refunds"
	applyBill     = "/v3/bill/tradebill"
	applyFundBill = "/v3/bill/fundflowbill"
)

type app struct {
	p *PayC
}

//统一下单提交数据
type UniformOrderData struct {
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

//订单金额
type Amount struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
	Refund   string `json:"refund,omitempty"`
}

//优惠功能
type Detail struct {
	CostPrice   string               `json:"cost_price,omitempty"`
	InvoiceId   string               `json:"invoice_id,omitempty"`
	GoodsDetail *map[int]GoodsDetail `json:"goods_detail,omitempty"`
}

//单品列表
type GoodsDetail struct {
	MerchantGoodsId  string `json:"merchant_goods_id,omitempty"`
	WechatpayGoodsId string `json:"wechatpay_goods_id,omitempty"`
	GoodsName        string `json:"goods_name,omitempty"`
	Quantity         string `json:"quantity,omitempty"`
	UnitPrice        string `json:"unit_price,omitempty"`
	RefundAmount     string `json:"refund_amount"`
	RefundQuantity   string `json:"refund_quantity"`
}

//场景信息
type Secene struct {
	PayerClientIp string    `json:"payer_client_ip"`
	DeviceId      string    `json:"device_id"`
	StoreInfo     StoreInfo `json:"store_info"`
}

//商户门店信息
type StoreInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	AeraName string `json:"aera_name"`
	Address  string `json:"address"`
}

type ApplyRefund struct {
	TransactionId string               `json:"transaction_id,omitempty"`
	OutTradeNo    string               `json:"out_trade_no,omitempty"`
	OutRefundNo   string               `json:"out_refund_no"`
	Reason        string               `json:"reason,omitempty"`
	NotifyUrl     string               `json:"notify_url,omitempty"`
	FundsAccount  string               `json:"funds_account,omitempty"`
	Amount        Amount               `json:"amount"`
	GoodsDetail   *map[int]GoodsDetail `json:"goods_detail,omitempty"`
}

//统一下单
func (a *app) UniformOrder(data UniformOrderData, timestamp string, nonce string) (Reply, error) {
	data.AppId = a.p.appId
	data.MchId = a.p.mchId
	if data.NotifyUrl == "" {
		data.NotifyUrl = a.p.notifyUrl
	}
	str, _ := json.Marshal(data)
	return a.p.send(post, uniformOrder, timestamp, nonce, str)
}

//只提交必交字段,若无其它硬性要求,可使用此方法构建基础内容
func (u *UniformOrderData) Default(total int, tradeName string, tradeNo string, notifyUrl ...string) {
	u.Amount.Total = total
	u.Description = tradeName
	u.OutTradeNo = tradeNo
	if len(notifyUrl) > 0 {
		u.NotifyUrl = notifyUrl[0]
	}
	u.Detail = nil
	u.SceneInfo = nil
}

//查询单号
func (a *app) QueryTradeNo(wxOrd string, nonce string) (Reply, error) {
	str := wxOrd + "?mchid=" + a.p.mchId
	return a.p.send(get, queryTradeNo+str, general.Date().Timestamp("s"), nonce, nil)
}

//关闭订单
func (a *app) CloseTrade(wxOrd string, nonce string) (Reply, error) {
	str := wxOrd + "/close"
	return a.p.send(post, queryTradeNo+str, general.Date().Timestamp("s"), nonce, []byte(`{"mchid":"`+a.p.mchId+`"}`))
}

//申请退款
func (a *app) ApplyRefund(ar *ApplyRefund, timestamp string, nonce string) (Reply, error) {
	if len(*ar.GoodsDetail) == 0 {
		ar.GoodsDetail = nil
	}
	str, _ := json.Marshal(ar)
	return a.p.send(post, applyRefund, timestamp, nonce, str)
}

//查询退款
func (a *app) QueryRefund(wxOrd string, nonce string) (Reply, error) {
	str := applyRefund + "/" + wxOrd
	return a.p.send(get, str, general.Date().Timestamp("s"), nonce, nil)
}

//申请交易账单
//  date为账单日期,必填
//  nonce为随即字符串
//  非必填options顺序为1-账单类型,2-压缩类型,顺序必须正确,多传的一律忽略
func (a *app) ApplyBill(date string, nonce string, options ...string) (Reply, error) {
	str := "?bill_date=" + date
	o := len(options)
	if o > 0 {
		str += "&bill_type=" + options[0]
		if o > 1 {
			str += "&tar_type=" + options[1]
		}
	}
	return a.p.send(get, applyBill+str, general.Date().Timestamp("s"), nonce, nil)
}

//申请资金账单
//  date为账单日期,必填
//  nonce为随即字符串
//  非必填options顺序为1-资金账户类型,2-压缩类型,顺序必须正确,多传的一律忽略
func (a *app) ApplyFundBill(date string, nonce string, options ...string) (Reply, error) {
	str := "?bill_date=" + date
	o := len(options)
	if o > 0 {
		str += "&account_type=" + options[0]
		if o > 1 {
			str += "&tar_type=" + options[1]
		}
	}
	return a.p.send(get, applyFundBill+str, general.Date().Timestamp("s"), nonce, nil)
}
