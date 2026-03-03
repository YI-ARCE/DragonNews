package timing

import (
	"errors"
	"sync"
	"time"
)

type Pool struct {
	pool map[string]*TimerFunc
	lock sync.RWMutex
}

type TimerFunc struct {
	//需要执行的方法
	Func func() bool
	//定时执行时间
	Timers time.Duration
	//当前定时器的工作状态
	Status bool
	//定时器实例
	ticker *time.Ticker
	//停止通道
	stopChan chan struct{}
}

var pools = Pool{
	pool: make(map[string]*TimerFunc),
	lock: sync.RWMutex{},
}

func init() {
	// 启动管理器
	go manager()
}

func manager() {
	// 定期清理已停止的定时器
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		pools.lock.Lock()
		for name, tf := range pools.pool {
			if !tf.Status && tf.ticker == nil {
				delete(pools.pool, name)
			}
		}
		pools.lock.Unlock()
	}
}

// Create 创建一个循环计时器,name为计时器命名,f为方法,times为每次计时的的时间.
// times可以使用time.Second * X来获取,单位是秒,纯数字则单位为纳秒
// 注意,方法一定要返回一个bool来确定是否需要继续执行,true为继续,false为结束
func Create(name string, f func() bool, times time.Duration) (*TimerFunc, error) {
	pools.lock.Lock()
	defer pools.lock.Unlock()

	_, flag := pools.pool[name]
	if flag {
		return &TimerFunc{}, errors.New("名称已被使用!")
	}

	tf := &TimerFunc{
		Func:     f,
		Timers:   times,
		Status:   true,
		stopChan: make(chan struct{}),
	}

	pools.pool[name] = tf
	return tf, nil
}

func Anonymous(f func() bool, times time.Duration) *TimerFunc {
	return &TimerFunc{
		Func:     f,
		Timers:   times,
		Status:   true,
		stopChan: make(chan struct{}),
	}
}

// Get 查询循环计时器,并返回该计时器实例
func Get(name string) (*TimerFunc, error) {
	pools.lock.RLock()
	defer pools.lock.RUnlock()

	_, flag := pools.pool[name]
	if !flag {
		return &TimerFunc{}, errors.New("定时器不存在!")
	}

	return pools.pool[name], nil
}

// Start 开启计时器
func (tf *TimerFunc) Start() {
	// 确保状态为true
	tf.Status = true

	// 创建新的ticker和停止通道
	tf.ticker = time.NewTicker(tf.Timers)
	tf.stopChan = make(chan struct{})

	go func() {
		defer func() {
			if tf.ticker != nil {
				tf.ticker.Stop()
				tf.ticker = nil
			}
			close(tf.stopChan)
		}()

		// 立即执行一次
		if !tf.Func() {
			tf.Status = false
			return
		}

		for {
			select {
			case <-tf.ticker.C:
				if !tf.Status {
					return
				}
				if !tf.Func() {
					tf.Status = false
					return
				}
			case <-tf.stopChan:
				return
			}
		}
	}()
}

// Stop 停止计时器
func (tf *TimerFunc) Stop() {
	tf.Status = false
	if tf.stopChan != nil {
		close(tf.stopChan)
	}
	if tf.ticker != nil {
		tf.ticker.Stop()
		tf.ticker = nil
	}
}

// Del 删除定时器
func Del(key string) {
	pools.lock.Lock()
	defer pools.lock.Unlock()

	if tf, ok := pools.pool[key]; ok {
		tf.Stop()
		delete(pools.pool, key)
	}
}
