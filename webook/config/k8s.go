//go:build k8s

package config

var Config = config{
	DB: MySQLConfig{
		DNS: "root:root@tcp(k8s-mysql:13306)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:16379",
	},
	Port: "8080",
}
