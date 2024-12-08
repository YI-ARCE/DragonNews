package sakuraPostPetal

import (
	"dragonNews/table"
)

const Alias = alias(`spp2`)

const AliasName = `spp2`

const Table = `sakura_post_petal`

const TableAlias = string(Table + " `" + Alias + "`")

const SpId = `sp_id` //

const SppNum = `spp_num` // 花瓣数量

type alias string

// Structure 帖子收获的花瓣数
type Structure struct {
	//
	SpId int `json:"sp_id,omitempty"`
	// 花瓣数量
	SppNum int `json:"spp_num,omitempty"`
}

func (a alias) Keys(s ...string) string {
	return table.String((*string)(&a), &s)
}

func AliasKeys(alias string, s ...string) string {
	return table.String(&alias, &s)
}

// SpId
func (a alias) SpId() string {
	return *(*string)(&a) + SpId
}

// SppNum 花瓣数量
func (a alias) SppNum() string {
	return *(*string)(&a) + SppNum
}
