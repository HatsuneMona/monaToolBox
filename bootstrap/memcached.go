package bootstrap

import (
	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
	"monaToolBox/global"
)

func InitializeMemcached() *memcache.Client {
	if len(global.Config.Memcached.ServerList) < 1 {

		return nil
	}

	client := memcache.New(global.Config.Memcached.ServerList...)

	if err := client.Ping(); err != nil {
		global.Log.Fatal("init memcached error! ping output:", zap.Error(err))
		return nil
	}

	return client
}
