package runner

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)


func CacheSetHash(db string, args redis.Args) error {
	c, err := redis.Dial("tcp", db)
	if err != nil {
		fmt.Println("[ERROR] Error writing to redis, catch it next time around")
		return err
	}
	defer c.Close()

	c.Do("HMSET", args...)

	return nil
}

func CacheSet(db string, key string, value string) error {
	c, err := redis.Dial("tcp", db)
	if err != nil {
		fmt.Println("[ERROR] Error writing to redis, catch it next time around")
		return err
	}
	defer c.Close()

	c.Do("SET", key, value)

	return nil
}

func CacheGetHash(db string, key string) (Devices, error)  {
	s := Devices{}
	c, err := redis.Dial("tcp", db)
	if err != nil {
		return s, err
	}
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return s, err
	}

	redis.ScanStruct(values, &s)

	return s, nil
}

func CacheGetHashKey(db string, args redis.Args) (Devices, error)  {
	s := Devices{}
	c, err := redis.Dial("tcp", db)
	if err != nil {
		return s, err
	}
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", args...))
	if err != nil {
		return s, err
	}

	redis.ScanStruct(values, &s)
	return s, nil
}
