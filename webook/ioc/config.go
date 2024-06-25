package ioc

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func InitLocalConfig() {
	fp := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()

	viper.SetConfigFile(*fp)
	// 监听配置文件变更
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func InitRemoteConfig() {
	viper.SetConfigType("yaml")
	if err := viper.AddRemoteProvider("etcd3", "http://localhost:12379", "E:/Git/webook"); err != nil {
		panic(err)
	}
	if err := viper.ReadRemoteConfig(); err != nil {
		panic(err)
	}
}
