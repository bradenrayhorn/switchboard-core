package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
	}

	viper.SetDefault("token_expiration", time.Hour*24)

	_ = viper.BindEnv("mongo_host", "MONGO_HOST")
	_ = viper.BindEnv("mongo_port", "MONGO_PORT")
	_ = viper.BindEnv("mongo_username", "MONGO_USERNAME")
	_ = viper.BindEnv("mongo_password", "MONGO_PASSWORD")
	_ = viper.BindEnv("mongo_database", "MONGO_DATABASE")
	_ = viper.BindEnv("jwt_secret", "JWT_SECRET")

	if len(viper.GetString("jwt_secret")) < 32 {
		log.Println("Attention! JWT secret may be insecure.")
	}
}
