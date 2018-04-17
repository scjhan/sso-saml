package redis

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
)

var gRedisProtocol string
var gRedisHost string

// InitRedis init the redis configuration
func InitRedis(protocol string, host string) {
	gRedisProtocol = protocol
	gRedisHost = host
}

// GetString return a string value by key
func GetString(key string) (string, error) {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	if err != nil {
		return "", fmt.Errorf("connect redis error, error = %s", err.Error())
	}

	defer redisCli.Close()

	result, err := redigo.String(redisCli.Do("GET", key))
	if err != nil {
		return "", nil
	}
	return result, nil
}

// GetInt return a int value by key
func GetInt(key string) (int, error) {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	defer redisCli.Close()
	if err != nil {
		return 0, fmt.Errorf("connect redis error, error = %s", err.Error())
	}

	return redigo.Int(redisCli.Do("GET", key))
}

// SetString set a string value by string key
func SetString(key string, value string, expire int) error {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	if err != nil {
		return fmt.Errorf("connect redis error, error = %s", err.Error())
	}

	defer redisCli.Close()
	_, err = redisCli.Do("SETEX", key, expire, value)

	return err
}

// SetInt set a int value by string key
func SetInt(key string, value int, expire int) error {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	if err != nil {
		return fmt.Errorf("connect redis error, error = %s", err.Error())
	}

	defer redisCli.Close()
	_, err = redisCli.Do("SETEX", key, expire, value)

	return err
}

// Delete remove a key
func Delete(key string) error {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	if err != nil {
		return fmt.Errorf("connect redis error, error = %s", err.Error())
	}
	defer redisCli.Close()
	_, err = redisCli.Do("DEL", key)

	return err
}

// SetAdd add a mem of key to a redis set
func SetAdd(key string, mem string) error {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	if err != nil {
		return fmt.Errorf("connect redis error, error = %s", err.Error())
	}
	defer redisCli.Close()
	_, err = redisCli.Do("SADD", key, mem)

	return err
}

// SMembers return smembers
func SMembers(key string) ([]string, error) {
	redisCli, err := redigo.Dial(gRedisProtocol, gRedisHost)
	if err != nil {
		return nil, fmt.Errorf("connect redis error, error = %s", err.Error())
	}
	defer redisCli.Close()

	rv, err := redisCli.Do("SMEMBERS", key)

	sa := []string{}
	for _, iv := range rv.([]interface{}) {
		sa = append(sa, string(iv.([]byte)))
	}

	return sa, err
}
