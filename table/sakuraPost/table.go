package sakuraPost

import (
	"dragonNews/table"
)

const Alias = alias(`sp1`)

const AliasName = `sp1`

const Table = `sakura_post`

const TableAlias = string(Table + " `" + Alias + "`")

const SpId = `sp_id` // 帖子ID

const SuId = `su_id` // 用户ID

const SpContent = `sp_content` // 帖子文本内容

const SpImages = `sp_images` // 图片内容

const SpStatus = `sp_status` // 帖子状态,1->公共,2->私有

const CreateTime = `create_time` // 发帖时间

type alias string

// Structure 帖子内容
type Structure struct {
	// 帖子ID
	SpId int `json:"sp_id,omitempty"`
	// 用户ID
	SuId int `json:"su_id,omitempty"`
	// 帖子文本内容
	SpContent string `json:"sp_content,omitempty"`
	// 图片内容
	SpImages string `json:"sp_images,omitempty"`
	// 帖子状态,1->公共,2->私有
	SpStatus int `json:"sp_status,omitempty"`
	// 发帖时间
	CreateTime int `json:"create_time,omitempty"`
}

func (a alias) Keys(s ...string) string {
	return table.String((*string)(&a), &s)
}

func AliasKeys(alias string, s ...string) string {
	return table.String(&alias, &s)
}

// SpId 帖子ID
func (a alias) SpId() string {
	return *(*string)(&a) + SpId
}

// SuId 用户ID
func (a alias) SuId() string {
	return *(*string)(&a) + SuId
}

// SpContent 帖子文本内容
func (a alias) SpContent() string {
	return *(*string)(&a) + SpContent
}

// SpImages 图片内容
func (a alias) SpImages() string {
	return *(*string)(&a) + SpImages
}

// SpStatus 帖子状态,1->公共,2->私有
func (a alias) SpStatus() string {
	return *(*string)(&a) + SpStatus
}

// CreateTime 发帖时间
func (a alias) CreateTime() string {
	return *(*string)(&a) + CreateTime
}
