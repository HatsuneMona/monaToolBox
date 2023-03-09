package global

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"monaToolBox/bootstrap/config"
)

var (
	Config      config.Config
	Log         *zap.Logger
	ConfigViper *viper.Viper
	DB          *gorm.DB
	Redis       *redis.Client
	MemCached   *memcache.Client
)
