package config

type config struct {
	DB    MySQLConfig
	Redis RedisConfig
}

type MySQLConfig struct {
	DNS string
}

type RedisConfig struct {
	Addr string
}
