package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Limiter LimiterConfig
}

type ServerConfig struct {
	Port int
	Env  string
}

type DBConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type LimiterConfig struct {
	Rps     int
	Burst   int
	Enabled bool
}

func loadPath(env string) string {
	if env == "development" {
		return "../config/dev-cfg.yml"
	} else {
		return "../config/dev-cfg.yml"
	}
}

func loadConfig(filePath, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filePath)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil

}

func parseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func GetConfig() *Config {
	fileP := loadPath("APP_ENV")
	v, err := loadConfig(fileP, "yml")
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := parseConfig(v)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
