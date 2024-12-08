package sakuraUser

import (
	"dragonNews/table"
)

const Alias = alias(`su4`)

const AliasName = `su4`

const Table = `sakura_user`

const TableAlias = string(Table + " `" + Alias + "`")

const SuId = `su_id` // 用户ID

const SuPhone = `su_phone` // 用户使用的手机号码

const SuName = `su_name` // 昵称

const SuLevel = `su_level` // 用户等级

const SuHead = `su_head` //

const SuOneKey = `su_one_key` // 独立密钥

const SuAllKey = `su_all_key` // 全局密钥

const SuStatus = `su_status` // 用户状态,default->正常,disable->封禁

const CreateTime = `create_time` // 首次注册时间[dn_parse:time]

type alias string

// Structure 用户表
type Structure struct {
	// 用户ID
	SuId int `json:"su_id,omitempty"`
	// 用户使用的手机号码
	SuPhone string `json:"su_phone,omitempty"`
	// 昵称
	SuName string `json:"su_name,omitempty"`
	// 用户等级
	SuLevel int `json:"su_level,omitempty"`
	//
	SuHead string `json:"su_head,omitempty"`
	// 独立密钥
	SuOneKey string `json:"su_one_key,omitempty"`
	// 全局密钥
	SuAllKey string `json:"su_all_key,omitempty"`
	// 用户状态,default->正常,disable->封禁
	SuStatus string `json:"su_status,omitempty"`
	// 首次注册时间[dn_parse:time]
	CreateTime int `json:"create_time,omitempty"`
}

func (a alias) Keys(s ...string) string {
	return table.String((*string)(&a), &s)
}

func AliasKeys(alias string, s ...string) string {
	return table.String(&alias, &s)
}

// SuId 用户ID
func (a alias) SuId() string {
	return *(*string)(&a) + SuId
}

// SuPhone 用户使用的手机号码
func (a alias) SuPhone() string {
	return *(*string)(&a) + SuPhone
}

// SuName 昵称
func (a alias) SuName() string {
	return *(*string)(&a) + SuName
}

// SuLevel 用户等级
func (a alias) SuLevel() string {
	return *(*string)(&a) + SuLevel
}

// SuHead
func (a alias) SuHead() string {
	return *(*string)(&a) + SuHead
}

// SuOneKey 独立密钥
func (a alias) SuOneKey() string {
	return *(*string)(&a) + SuOneKey
}

// SuAllKey 全局密钥
func (a alias) SuAllKey() string {
	return *(*string)(&a) + SuAllKey
}

// SuStatus 用户状态,default->正常,disable->封禁
func (a alias) SuStatus() string {
	return *(*string)(&a) + SuStatus
}

// CreateTime 首次注册时间[dn_parse:time]
func (a alias) CreateTime() string {
	return *(*string)(&a) + CreateTime
}
