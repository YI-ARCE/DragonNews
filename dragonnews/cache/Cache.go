package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	. "yiarce/dragonnews/cache/driver"
)

type Redis struct {
	redis.Conn
}

//types 为连接类型,例如tcp
//host为连接地址,可以是IP,也可以是域名地址
func Init(types string, host string, options ...redis.DialOption) (Redis, error) {
	Conn, err := redis.Dial(types, host, options...)
	if err != nil {
		fmt.Println(err)
		return Redis{}, err
	} else {
		return Redis{Conn}, nil
	}
}

//获取数据
func (c Redis) Get(column string, ps interface{}) error {
	reply, err := c.Do("GET", column)
	if err != nil {
		return err
	}
	if reply == nil {
		return nil
	}
	if data, falg := reply.([]byte); falg {
		err = UnSerizlize(data, ps)
		if err != nil {
			return err
		}
	}
	return nil
}

//写入数据
func (c Redis) Set(column string, data interface{}) error {
	serialize, err := Serialize(data)
	if err != nil {
		return err
	}
	_, err = c.Do("SET", column, serialize)
	if err != nil {
		return err
	}
	return nil
}

//删除数据
func (c Redis) Del(column ...interface{}) error {
	_, err := c.Do("DEL", column...)
	if err != nil {
		return err
	}
	return nil
}
