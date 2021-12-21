package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBSource       string `mapstructure:"DB_SOURCE"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	ImgurClientID  string `mapstructure:"IMGUR_CLIENT_ID"`
	ImgurUploadURL string `mapstructure:"IMGUR_UPLOAD_URL"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.SetDefault("ServerAddress", "0.0.0.0")
	viper.SetDefault("ServerPort", "8080")

	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file

	// Handle errors reading the config file
	if err != nil {
		err = fmt.Errorf("reading config file: %w ", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		err = fmt.Errorf("unmarshalling config: %w ", err)
		return
	}

	return
}
