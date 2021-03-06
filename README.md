# Go Locks

- map lock
- redis lock
- mysql lock

## Usage

```go
//  Map Lock
mapLock := NewMapLock("lock1")
err := mapLock.TryLock()
err = mapLock.Unlock()

// Redis Lock
InitRedisLock(redisClient)
redisLock := NewRedisLock("lock2", time.Second)
err = redisLock.TryLock()
err = redisLock.Unlock()

// MySQL Lock
InitMysqlLock(db, "test", 5*time.Second)
mysqlLock := NewMysqlLock("lock3", time.Second)
err = mysqlLock.TryLock()
err = mysqlLock.Unlock()

// Upgrade to Spin Lock
mapSpinLock := NewSpinLock(NewMapLock(lockKey), spinTries, spinInterval)
err = mapSpinLock.Lock()
err = mapSpinLock.Unlock()

```

### Dependency Injection

```go
type usecase struct {
	lockFactory ExpiryLockFactory
}

func (u *usecase) CheckFrequntSubmit(key string, expiry time.Duration) (ok bool) {
	lock := u.lockFactory.NewLock(key, expiry)
	if err := lock.TryLock(); err != nil {
		return false
	}

	// let the lock expire, and unlock automatically
	return true
}

func TestUsecase_CheckFrequentSubmit(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	key := "book:1:submit"
	expiry := time.Second

	lock := mock.NewMockTryLocker(ctl)
	lock.EXPECT().TryLock().Return(nil)
	lockFactory := mock.NewMockExpiryLockFactory(ctl)
	lockFactory.EXPECT().NewLock(key, expiry).Return(lock)

	bookUsecase := usecase{lockFactory}
	ok := bookUsecase.CheckFrequntSubmit(key, expiry)
	assert.True(t, ok)
}
```