package sakuraUserProperty

import (
	"dragonNews/table"
)

const Alias = alias(`sup5`)

const AliasName = `sup5`

const Table = `sakura_user_property`

const TableAlias = string(Table + " `" + Alias + "`")

const SupId = `sup_id` // 拥有ID

const SupPetal = `sup_petal` // 花瓣

const SupFlowers = `sup_flowers` // 花朵数量

type alias string

// Structure 所有用户的固定拥有物品
type Structure struct {
	// 拥有ID
	SupId int `json:"sup_id,omitempty"`
	// 花瓣
	SupPetal int `json:"sup_petal,omitempty"`
	// 花朵数量
	SupFlowers int `json:"sup_flowers,omitempty"`
}

func (a alias) Keys(s ...string) string {
	return table.String((*string)(&a), &s)
}

func AliasKeys(alias string, s ...string) string {
	return table.String(&alias, &s)
}

// SupId 拥有ID
func (a alias) SupId() string {
	return *(*string)(&a) + SupId
}

// SupPetal 花瓣
func (a alias) SupPetal() string {
	return *(*string)(&a) + SupPetal
}

// SupFlowers 花朵数量
func (a alias) SupFlowers() string {
	return *(*string)(&a) + SupFlowers
}
