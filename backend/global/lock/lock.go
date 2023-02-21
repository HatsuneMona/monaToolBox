// Package lock
// @Author HatsuneMona 2022/11/11 0:14
package lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"monaToolBox/global"
	"monaToolBox/utils"
	"time"
)

type LockInterface interface {
	Lock() bool
	Block(seconds int64) bool
	Release() bool
	ForceRelease()
}

type lock struct {
	context context.Context
	name    string // 锁名称
	owner   string // 锁标识
	seconds int64  // 有效期
}

const LOCK_PREFIX = "LOCK_"

// 释放锁 Lua 脚本，防止任何客户端都能解锁
// KEYS[1]: lock.name
// ARGV[1]: lock.owner
const releaseLockLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`

// GetLock 生成锁
func GetLock(name string, seconds int64) LockInterface {
	return &lock{
		context.Background(),
		LOCK_PREFIX + name,
		utils.RandString(16),
		seconds,
	}
}

// Lock 获取锁
// 若锁(key)不存在，则添加key-value，并返回 true，代表已获取到锁
// 若锁(key)已存在，则不进行任何操作，并返回 false，代表未获取到锁
func (l *lock) Lock() bool {
	// SetNX：SET if Not eXists
	return global.Redis.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

// Block 阻塞1秒后，尝试重新获取锁，持续seconds秒。
func (l *lock) Block(seconds int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Lock() {
			time.Sleep(time.Duration(1) * time.Second)
			if time.Now().Unix()-seconds >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

// Release 释放锁
func (l *lock) Release() bool {
	luaScript := redis.NewScript(releaseLockLuaScript)

	result := luaScript.Run(
		l.context, global.Redis,
		[]string{l.name}, // []string{l.name}   KEYS[1]
		l.owner,          // l.owner            ARGV[1]
	).Val().(int64)
	return result != 0
}

// ForceRelease 强制释放锁
func (l *lock) ForceRelease() {
	global.Redis.Del(l.context, l.name).Val()
}
