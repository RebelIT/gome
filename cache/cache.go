package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Devices struct{
	Device		string `json:"device"`
	Name		string `json:"name"`
	Addr 		string `json:"address"`
	NetPort		string `json:"port"`
	Id			string `json:"id"`
	Key 		string `json:"key"`
}

type Inputs struct{
	Database	string `json:"database"`
	Devices		[]Devices
}

type Status struct{
	Device	string `json:"device"`
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

func SetHash(db string, args redis.Args) error {
	c, err := redis.Dial("tcp", db)
	if err != nil {
		fmt.Println("[ERROR] Error writing to redis, catch it next time around")
		return err
	}
	defer c.Close()

	c.Do("HMSET", args...)

	return nil
}

func Set(db string, key string, value string) error {
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

func GetStatus(db string, key string) (Status, error)  {
	s := Status{}
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

func GetHashKey(db string, args redis.Args) (Devices, error)  {
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
