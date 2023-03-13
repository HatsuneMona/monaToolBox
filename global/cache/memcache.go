package cache

import (
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"monaToolBox/global"
	"monaToolBox/utils"
)

type cacheMc struct {
	c *memcache.Client

	prefix string
}

type IMemCache interface {
	buildKey(string) string
	Set(string, any, int32) error
	// AddMulti(map[string]any, int) error

	Get(string) ([]byte, error)
	GetMulti([]string) (map[string][]byte, error)

	Delete(string) error
	// DeleteMulti([]string) error
}

// NewMemcachedClient 打开一个新的memcached对话
func NewMemcachedClient(prefix string) *cacheMc {
	if prefix == "" {
		global.Log.Error("create memcached without prefix")
		return nil
	}
	mc := cacheMc{
		prefix: utils.BuildString(prefix, "_"),
		c:      global.MemCached,
	}

	return &mc
}

// Set 添加一个缓存
func (mc *cacheMc) Set(key string, item any, expiration int32) error {
	if key == "" || item == nil || expiration < 1 {
		return errors.New("input error")
	}

	var value []byte
	if v, ok := item.(string); ok {
		value = []byte(v)
	} else if v, ok := item.([]byte); ok {
		value = v
	} else {
		v, err := json.Marshal(item)
		if err != nil {
			return err
		}
		value = v
	}

	mcItem := memcache.Item{
		Key:        mc.buildKey(key),
		Value:      value,
		Expiration: expiration,
	}

	return mc.c.Set(&mcItem)
}

// func (mc *cacheMc) AddMulti(m map[string]any) error {
// 	// TOD implement me
// 	panic("implement me")
// }

// Get 给定 key 获取缓存结果
func (mc *cacheMc) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, errors.New("input error")
	}

	get, err := mc.c.Get(mc.buildKey(key))
	if err != nil {
		return nil, err
	}

	return get.Value, nil
}

// GetMulti 给一批key，获取一堆结果
func (mc *cacheMc) GetMulti(keys []string) (map[string][]byte, error) {
	if len(keys) < 1 {
		return nil, errors.New("input error")
	}
	realKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		realKeys = append(realKeys, mc.buildKey(key))
	}

	gets, err := mc.c.GetMulti(realKeys)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for key, i := range gets {
		result[key[len(mc.prefix):]] = i.Value // 把前缀去掉
	}
	return result, nil
}

// Delete 删除一个key
func (mc *cacheMc) Delete(key string) error {
	if key == "" {
		return errors.New("input error")
	}
	return mc.c.Delete(mc.buildKey(key))
}

// func (mc *cacheMc) DeleteMulti(strings []string) error {
// 	// TOD implement me
// 	panic("implement me")
// }

func (mc *cacheMc) buildKey(key string) string {
	return utils.BuildString(mc.prefix, key)
}
