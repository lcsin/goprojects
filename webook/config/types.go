package config

type config struct {
	DB     MySQLConfig
	Redis  RedisConfig
	JWTKey string
	Port   string
}

type MySQLConfig struct {
	DNS string
}

type RedisConfig struct {
	Addr   string
	Passwd string
}
