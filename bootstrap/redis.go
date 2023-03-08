package bootstrap

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"monaToolBox/global"
	"strconv"
)

// InitializeRedis 初始化Redis
func InitializeRedis() *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:     global.Config.Redis.Host + ":" + strconv.Itoa(global.Config.Redis.Port),
			Password: global.Config.Redis.Password,
			DB:       global.Config.Redis.DB,
		},
	)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
		return nil
	}
	return client
}
