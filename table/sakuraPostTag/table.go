package sakuraPostTag

import (
	"dragonNews/table"
)

const Alias = alias(`spt3`)

const AliasName = `spt3`

const Table = `sakura_post_tag`

const TableAlias = string(Table + " `" + Alias + "`")

const SpId = `sp_id` // 帖子ID

const SptContent = `spt_content` // tag集合文本

const SptTagOne = `spt_tag_one` // tag1

const SptTagTwo = `spt_tag_two` // tag2

const SptTagThree = `spt_tag_three` // tag3

const SptTagFour = `spt_tag_four` // tag4

const SptTagFive = `spt_tag_five` // tag5

type alias string

// Structure 帖子TAG中间表
type Structure struct {
	// 帖子ID
	SpId int `json:"sp_id,omitempty"`
	// tag集合文本
	SptContent string `json:"spt_content,omitempty"`
	// tag1
	SptTagOne string `json:"spt_tag_one,omitempty"`
	// tag2
	SptTagTwo string `json:"spt_tag_two,omitempty"`
	// tag3
	SptTagThree string `json:"spt_tag_three,omitempty"`
	// tag4
	SptTagFour string `json:"spt_tag_four,omitempty"`
	// tag5
	SptTagFive string `json:"spt_tag_five,omitempty"`
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

// SptContent tag集合文本
func (a alias) SptContent() string {
	return *(*string)(&a) + SptContent
}

// SptTagOne tag1
func (a alias) SptTagOne() string {
	return *(*string)(&a) + SptTagOne
}

// SptTagTwo tag2
func (a alias) SptTagTwo() string {
	return *(*string)(&a) + SptTagTwo
}

// SptTagThree tag3
func (a alias) SptTagThree() string {
	return *(*string)(&a) + SptTagThree
}

// SptTagFour tag4
func (a alias) SptTagFour() string {
	return *(*string)(&a) + SptTagFour
}

// SptTagFive tag5
func (a alias) SptTagFive() string {
	return *(*string)(&a) + SptTagFive
}
