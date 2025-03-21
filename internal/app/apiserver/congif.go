package apiserver

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	RedisURL    string `toml:"redis"`
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
