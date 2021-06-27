package apiserver

type Config struct {
	SrvHost    string `toml:"srv_host"`
	SrvPort    string `toml:"srv_port"`
	DBHost     string `toml:"db_host"`
	DBPort     string `toml:"db_port"`
	DBUser     string `toml:"db_user"`
	DBPassword string `toml:"db_password"`
	DBName     string `toml:"db_name"`
}

func NewConfig() *Config {
	return &Config{}
}
