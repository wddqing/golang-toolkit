package xredis

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/wddqing/golang-toolkit/share"
)

// RedisConfig RedisConfig
type RedisConfig struct {
	Address  string `json:"addr"`
	Password string `json:"password"`
	DBNum    int    `json:"db"`
}

type Redis struct {
	*redis.Pool
}

func NewRedis(conf RedisConfig, options ...Option) *Redis {
	opts := defaultOptions
	for _, option := range options {
		option(&opts)
	}

	pool := &redis.Pool{
		MaxActive:   opts.maxActive,
		MaxIdle:     opts.maxIdle,
		IdleTimeout: opts.idleTimeout,
		Wait:        opts.wait,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", conf.Address, redis.DialPassword(conf.Password), redis.DialDatabase(conf.DBNum))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &Redis{
		Pool: pool,
	}
}

func NewRedisByPool(pool *redis.Pool) *Redis {
	return &Redis{
		Pool: pool,
	}
}

func (r *Redis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := r.Get()
	reply, err = conn.Do(commandName, args...)
	conn.Close()
	return
}

func (r *Redis) IncUtilToNight(key string, v int) (int, error) {
	today := share.GetTodayBeginAt()
	now := time.Now()
	return r.IncWithInitialTime(key, v, int(86400-(now.Unix()-today.Unix())))
}

//设置目标key增加v,并且当key不存在的时候生存时间为t秒,返回增加后的值
func (r *Redis) IncWithInitialTime(key string, expired int, v int) (int, error) {
	ret, err := redis.Int(r.Do("INCRBY", key, v))
	if err != nil {
		return ret, err
	}
	if ret == v { //说明key不存在
		_, err = r.Do("EXPIRE", key, expired)
		return ret, err
	}
	return ret, nil
}

func (r *Redis) Del(key string) error {
	_, err := r.Do("DEL", key)
	return err
}

// 上锁, expired最大有效时间
func (r *Redis) lock(key string, expired int) error {
	ret, err := redis.Int(r.Do("INCRBY", key, 1))
	if err != nil {
		return err
	}
	switch ret {
	case 1:
		// 设置最大有效时间
		_, err = r.Do("EXPIRE", key, expired)
		if err != nil {
			return err
		}
	case 0:
		return errors.New("Redis 上锁失败")
	default:
		return errors.New("此操作正在执行，请稍后再试")
	}

	return nil
}

// f()函数操作 防止并发上锁，默认30秒
func (r *Redis) DoWithLock(f func() error, key string) error {
	// 防止并发加锁
	err := r.lock(key, 30)
	if err != nil {
		return err
	}
	// 解锁
	defer r.Del(key)
	err = f()
	if err != nil {
		return err
	}
	return nil
}

// f()函数操作 防止并发上锁和限制频率，expiredAt秒内操作一次
func (r *Redis) DoInExpiredAt(f func() error, key string, expiredAt int) error {
	// 防止并发加锁
	err := r.lock(key, expiredAt)
	if err != nil {
		return err
	}
	err = f()
	if err != nil {
		// 失败解锁
		r.Del(key)
		return err
	}
	return nil
}

func (r *Redis) TryLock(key string, expired int) bool {
	if err := r.lock(key, expired); err != nil {
		return false
	}
	return true
}
