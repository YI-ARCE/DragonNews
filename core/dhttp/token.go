package dhttp

import (
	"math/rand"
	"time"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())
}
