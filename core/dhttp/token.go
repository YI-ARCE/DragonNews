package dhttp

import (
	"math/rand"
	"time"
)

func init() {
	// 设置随机数种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
