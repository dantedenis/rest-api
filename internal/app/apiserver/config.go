package apiserver

type Config struct {
	BinAddr  string `yaml:"bin_addr"`
	LogLevel string `yaml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		BinAddr:  ":8080",
		LogLevel: "debug",
	}
}
