package model

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Hostname string `yaml:"hostname"`
		CORS    struct {
			Enabled      bool     `yaml:"enabled"`
			AllowOrigins []string `yaml:"allowOrigins"`
		} `yaml:"cors"`
	} `yaml:"server"`
	Database struct {
		DSN string `json:"dsn"`
	} `yaml:"database"`
	Email struct {
		SMTP struct {
			Hostname string `yaml:"hostname"`
			Port int `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"smtp"`
	} `yaml:"email"`
}
