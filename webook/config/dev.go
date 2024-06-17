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
	JWTKey: "fsAck3=%n*&*6XxbCd5ksXGjLHZT2fXc",
	Port:   "8080",
}
