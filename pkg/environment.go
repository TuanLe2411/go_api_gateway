package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Routes []Route `mapstructure:"routes"`
}

type Route struct {
	Name    string `mapstructure:"name"`
	Context string `mapstructure:"context"`
	Target  string `mapstructure:"target"`
}

func loadRoutes(env string) (*Config, error) {
	viper.SetConfigName(env)
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	case "development":
		envFile = ".env.development"
	default:
		log.Fatalf("ENV không hợp lệ: %s. Chỉ hỗ trợ 'development' hoặc 'production'", env)
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Lỗi khi load file %s: %v", envFile, err)
	}

	return loadRoutes(env)
}
