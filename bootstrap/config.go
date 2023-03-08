package bootstrap

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"monaToolBox/global"
	"os"
)

func InitConfig() *viper.Viper {

	homeDir, err := homedir.Dir()
	if err != nil {
		// TODO log
		panic(err)
	}
	configPath := homeDir + "/.monaPanel/"

	appPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	v := viper.New()
	v.AddConfigPath(configPath)
	v.AddConfigPath(appPath)
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config faild : %s \n", err))
	}

	if err := v.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unmarshal config error, %s \n", err)
	}

	return v
}
