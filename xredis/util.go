package xredis

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

// 反序列化
func GetAndJsonUnmarshal(r *Redis, key string, obj interface{}) error {
	data, err := redis.Bytes(r.Do("GET", key))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, obj)
}

// 序列化
func SetExAndJsonMarshalHelper(r *Redis, key string, expire int32, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = r.Do("SETEX", key, expire, data)
	return err
}
