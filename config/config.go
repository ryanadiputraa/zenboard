package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Database Database
	Oauth    Oauth
	JWT      JWT
}

type Server struct {
	Env  string
	Port string
}

type Database struct {
	Driver string
	DSN    string
}

type Oauth struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	RedirectURL  string
	State        string
}

type JWT struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshSecret    string
	RefreshExpiresIn time.Duration
}

func LoadConfig(filename string) (config *Config, err error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err = v.ReadInConfig(); err != nil {
		return
	}

	err = v.Unmarshal(&config)
	return
}
