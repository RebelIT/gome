package database

import (
	"fmt"
	"github.com/rebelit/gome/util/config"
	db "github.com/rebelit/gome_badger"
)

func Get(name string) (data string, error error) {
	conn, err := db.Open(config.App.DbPath)
	if err != nil {
		return "", fmt.Errorf("open database: %s", err)
	}
	defer conn.Close()

	value, err := conn.Get(name)
	if err != nil {
		return "", fmt.Errorf("get database value: %s", err)
	}
	return value, nil
}

func Add(key, value string) error {
	conn, err := db.Open(config.App.DbPath)
	if err != nil {
		return fmt.Errorf("open database: %s", err)
	}
	defer conn.Close()

	if err := conn.Set(key, value); err != nil {
		return fmt.Errorf("add database value: %s", err)
	}
	return nil
}

func Del(key string) error {
	conn, err := db.Open(config.App.DbPath)
	if err != nil {
		return fmt.Errorf("open database: %s", err)
	}
	defer conn.Close()

	if err := conn.Delete(key); err != nil {
		return fmt.Errorf("del database key: %s", err)
	}
	return nil
}

func GetAll() (data []string, error error) {
	conn, err := db.Open(config.App.DbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %s", err)
	}
	defer conn.Close()

	value, err := conn.GetAllKeys()
	if err != nil {
		return nil, fmt.Errorf("get all database values: %s", err)
	}
	return value, nil
}
