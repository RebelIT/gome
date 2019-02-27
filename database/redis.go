package database

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"log"
)

// *****************************************************************
// Redis functions
func DbConnect()(redis.Conn, error){
	s, err := common.GetSecrets()
	if err != nil{
		log.Printf("[ERROR] slack alert: %s\n", err)
		return nil, err
	}

	db := s.Database
	conn, err := redis.Dial("tcp", db)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DbHashSet(key string, data interface{} ) error{
	//equivalent to a redis HMSET
	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()
	if _, err := c.Do("HMSET", redis.Args{key}.AddFlat(data)...); err != nil{
		return err
	}

	return nil
}

func DbHashGet(key string)(values []interface{}, err error){
	c, err := DbConnect()
	if err != nil{
		return nil, err
	}
	defer c.Close()

	resp, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func DbGetKeys(key string)(keys []string, err error){
	c, err := DbConnect()
	if err != nil{
		return nil, err
	}
	defer c.Close()

	resp, err := redis.Strings(c.Do("KEYS", key))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DbSet(key string, value []byte) error{
	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()
	if _, err := c.Do("SET", key, string(value)); err != nil{
		return err
	}
	return nil
}

func DbGet(key string) (values string, err error){
	c, err := DbConnect()
	if err != nil{
		return "", err
	}
	defer c.Close()
	value, err := redis.String(c.Do("GET", key))
	if err != nil{
		if err.Error() == "redigo: nil returned"{
			//This is fine, redis connection OK, just no data returned
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func DbDel(key string) error{
	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()

	if _, err := c.Do("DEL", key, "*"); err != nil{
		return err
	}

	return nil
}