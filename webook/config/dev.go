//go:build !k8s

package config

var Cfg = config{
	DB: MySQLConfig{
		//DNS: "root:root@tcp(localhost:31306)/webook",
		DNS: "root:root@tcp(112.124.62.35:3306)/webook",
	},
	Redis: RedisConfig{
		//Addr:   "localhost:31379",
		Addr:   "112.124.62.35:16379",
		Passwd: "root123",
	},
	JWTKey: "fsAck3=%n*&*6XxbCd5ksXGjLHZT2fXc",
	Port:   "8080",
}
