package config

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	HTTPPort    string        `mapstructure:"http_port" validate:"required"`
	GRPCPort    string        `mapstructure:"grpc_port" validate:"required"`
	Timeout     time.Duration `mapstructure:"timeout" validate:"required"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout" validate:"required"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host" validate:"omitempty,hostname"`
	Port     int    `mapstructure:"port" validate:"omitempty,min=1024,max=65535"`
	User     string `mapstructure:"user" validate:"omitempty,min=3"`
	Password string `mapstructure:"password" validate:"omitempty,min=6"`
	DBName   string `mapstructure:"dbname" validate:"omitempty,min=3"`
}

type StorageConfig struct {
	Type     string         `mapstructure:"type" validate:"required,oneof=memory postgres"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type LogConfig struct {
	Level string `mapstructure:"level" validate:"required,oneof=local prod"`
}

type Config struct {
	Server  ServerConfig  `mapstructure:"server" validate:"required"`
	Storage StorageConfig `mapstructure:"storage" validate:"required"`
	Log     LogConfig     `mapstructure:"log" validate:"required"`
}

func MustLoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config, %v", err)
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("error validating config, %v", err)
	}

	return &cfg
}
