package help

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

var Redis *RedisModel

func init() {
	var err error
	Redis, err = NewRedis()
	Error(err)
}

func NewRedis() (rc *RedisModel, err error) {
	host := beego.AppConfig.String("redis.host")
	port := beego.AppConfig.String("redis.port")
	prefix := beego.AppConfig.String("redis.key.prefix")
	if host == "" || port == "" {
		return nil, errors.New("配置错误")
	}

	rc = new(RedisModel)

	dialFunc := func() (c redis.Conn, err error) {
		//        c, err = redis.Dial("tcp", rc.conninfo)
		c, err = redis.Dial("tcp", host+":"+port)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}

	rc.prefix = prefix

	return rc, nil
}

type RedisModel struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	prefix   string
}

func (rc *RedisModel) Set(key, val string) error {
	c := rc.p.Get()
	defer c.Close()
	prefix := rc.prefix

	_, err := c.Do("SET", prefix+key, val)
	Error(err)

	return err
}

func (rc *RedisModel) Get(key string) string {
	c := rc.p.Get()
	defer c.Close()
	prefix := rc.prefix
	val, err := redis.String(c.Do("GET", prefix+key))

	Error(err)

	return val
}

func (rc *RedisModel) Lpush(key, val string) error {
	c := rc.p.Get()
	defer c.Close()
	prefix := rc.prefix
	_, err := c.Do("lpush", prefix+key, val)

	Error(err)
	return err
}

func (rc *RedisModel) Rpop(key string) string {
	c := rc.p.Get()
	defer c.Close()
	prefix := rc.prefix
	val, err := redis.String(c.Do("rpop", prefix+key))

	Error(err)
	return val

}