package timing

import (
	"errors"
	"time"
)

type Pool struct {
	pool map[string]*TimerFunc
}

type TimerFunc struct {
	//需要执行的方法
	Func func() bool
	//定时执行时间
	Timers time.Duration
	//当前定时器的工作状态
	Status bool
}

var pools = Pool{
	map[string]*TimerFunc{},
}

func manager() {

}

// Create 创建一个循环计时器,name为计时器命名,f为方法,times为每次计时的的时间.
// times可以使用time.Second * X来获取,单位是秒,纯数字则单位为纳秒
// 注意,方法一定要返回一个bool来确定是否需要继续执行,true为继续,false为结束
func Create(name string, f func() bool, times time.Duration) (*TimerFunc, error) {
	_, flag := pools.pool[name]
	if flag {
		return &TimerFunc{}, errors.New("名称已被使用!")
	}
	pools.pool[name] = &TimerFunc{
		Func:   f,
		Timers: times,
		Status: true,
	}
	return pools.pool[name], nil
}

func Anonymous(f func() bool, times time.Duration) *TimerFunc {
	return &TimerFunc{
		Func:   f,
		Timers: times,
		Status: true,
	}
}

// Get 查询循环计时器,并返回该计时器实例
func Get(name string) (*TimerFunc, error) {
	_, flag := pools.pool[name]
	if !flag {
		return &TimerFunc{}, errors.New("定时器不存在!")
	}
	return pools.pool[name], nil
}

// Start 开启计时器
func (tf *TimerFunc) Start() {
	tf.Status = true
	t := time.NewTicker(tf.Timers)
	go func() {
		defer func() {
			t.Stop()
		}()
		if !tf.Func() {
			return
		}
		for range t.C {
			if !tf.Status {
				return
			}
			if !tf.Func() {
				return
			}
		}
	}()
}

// Stop 停止计时器
func (tf *TimerFunc) Stop() {
	tf.Status = false
}

func Del(key string) {
	delete(pools.pool, key)
}
