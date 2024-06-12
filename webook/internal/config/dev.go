//go:build !k8s

package config

var Config = config{
	DB: MySQLConfig{
		DNS: "root:root@tcp(localhost:31306)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:31379",
	},
}
