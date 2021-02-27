package weChatPay

import (
	"encoding/json"
	"yiarce/dragonnews/general"
)

const (
	queryMchTradeNo = "/v3/pay/transactions/out-trade-no/"
	queryWcTradeNo  = "/v3/pay/transactions/id/"
	applyRefund     = "/v3/refund/domestic/refunds"
	applyBill       = "/v3/bill/tradebill"
	applyFundBill   = "/v3/bill/fundflowbill"
)

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

type generals struct {
	p *PayC
}

//查询商户订单号
func (g *generals) QueryMchTradeNo(wxOrd string, nonce string) (Reply, error) {
	str := wxOrd + "?mchid=" + g.p.mchId
	return g.p.send(get, queryMchTradeNo+str, general.Date().Timestamp("s"), nonce, nil)
}

//查询微信支付单号
func (g *generals) QueryWeChatTradeNo(wxOrd string, nonce string) (Reply, error) {
	str := wxOrd + "?mchid=" + g.p.mchId
	return g.p.send(get, queryWcTradeNo+str, general.Date().Timestamp("s"), nonce, nil)
}

//关闭订单
func (g *generals) CloseTrade(wxOrd string, nonce string) (Reply, error) {
	str := wxOrd + "/close"
	return g.p.send(post, queryMchTradeNo+str, general.Date().Timestamp("s"), nonce, []byte(`{"mchid":"`+g.p.mchId+`"}`))
}

//申请退款
func (g *generals) ApplyRefund(ar *ApplyRefund, timestamp string, nonce string) (Reply, error) {
	if len(*ar.GoodsDetail) == 0 {
		ar.GoodsDetail = nil
	}
	str, _ := json.Marshal(ar)
	return g.p.send(post, applyRefund, timestamp, nonce, str)
}

//查询退款
func (g *generals) QueryRefund(wxOrd string, nonce string) (Reply, error) {
	str := applyRefund + "/" + wxOrd
	return g.p.send(get, str, general.Date().Timestamp("s"), nonce, nil)
}

//申请交易账单
//  date为账单日期,必填
//  nonce为随即字符串
//  非必填options顺序为1-账单类型,2-压缩类型,顺序必须正确,多传的一律忽略
func (g *generals) ApplyBill(date string, nonce string, options ...string) (Reply, error) {
	str := "?bill_date=" + date
	o := len(options)
	if o > 0 {
		str += "&bill_type=" + options[0]
		if o > 1 {
			str += "&tar_type=" + options[1]
		}
	}
	return g.p.send(get, applyBill+str, general.Date().Timestamp("s"), nonce, nil)
}

//申请资金账单
//  date为账单日期,必填
//  nonce为随即字符串
//  非必填options顺序为1-资金账户类型,2-压缩类型,顺序必须正确,多传的一律忽略
func (g *generals) ApplyFundBill(date string, nonce string, options ...string) (Reply, error) {
	str := "?bill_date=" + date
	o := len(options)
	if o > 0 {
		str += "&account_type=" + options[0]
		if o > 1 {
			str += "&tar_type=" + options[1]
		}
	}
	return g.p.send(get, applyFundBill+str, general.Date().Timestamp("s"), nonce, nil)
}
