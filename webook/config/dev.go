//go:build !k8s

package config

var Cfg = config{
	DB: MySQLConfig{
		DNS: "root:root@tcp(localhost:31306)/webook",
	},
	Redis: RedisConfig{
		Addr:   "localhost:31379",
		Passwd: "",
	},
	Port: "8080",
}
