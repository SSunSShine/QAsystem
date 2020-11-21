package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Configuration struct {
	DB struct {
		Driver string `json:"driver"`
		Addr   string `json:"addr"`
	} `json:"db"`
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redis"`
	Address string `json:"address"`
	JwtKey  string `json:"jwtKey"`
}

var conf *Configuration

// Config 加载配置信息
func Config() *Configuration {
	if conf != nil {
		return conf
	}

	viper.SetConfigName("configuration")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config file error:", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Println("config file error:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration information loaded successfully!", conf)

	return conf
}
