package golocks

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var (
	testRedisHost = env("TEST_REDIS_HOST", "127.0.0.1")
	testRedisPORT = env("TEST_REDIS_PORT", "6379")
	testRedisPWD  = env("TEST_REDIS_PWD", "")
)

func TestRedisLock_Lock(t *testing.T) {
	client := getRedisClient(t, testRedisHost, testRedisPORT, testRedisPWD)
	InitRedisLock(client)
	lock := NewRedisLock("test", time.Second)

	err := lock.TryLock()
	assert.Nil(t, err)
	err = lock.Unlock()
	assert.Nil(t, err)
}

func TestRedisLock_Expired(t *testing.T) {
	client := getRedisClient(t, testRedisHost, testRedisPORT, testRedisPWD)
	InitRedisLock(client)
	lock := NewRedisLock("expiry", 500*time.Millisecond)

	err := lock.TryLock()
	assert.Nil(t, err)
	time.Sleep(600 * time.Millisecond)

	err = lock.Unlock()
	assert.NotNil(t, err)
	err = lock.TryLock()
	assert.Nil(t, err)
}

func getRedisClient(t *testing.T, host, port, pwd string) *redis.Client {
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	if err := client.Ping().Err(); err != nil {
		t.Fatal(err)
	}

	client.FlushAll()
	return client
}
