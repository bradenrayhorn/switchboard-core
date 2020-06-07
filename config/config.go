package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
	}

	_ = viper.BindEnv("mongo_host", "MONGO_HOST")
	_ = viper.BindEnv("mongo_port", "MONGO_PORT")
	_ = viper.BindEnv("jwt_secret", "JWT_SECRET")

	if len(viper.GetString("jwt_secret")) < 32 {
		log.Println("Attention! JWT secret may be insecure.")
	}
}
