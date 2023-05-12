package configs

import "github.com/spf13/viper"

type Config struct {
	Port          string `mapstruct:"PORT"`
	DBUrl         string `mpastruct:"DB_URL"`
	ProductSvcUrl string `mapstruct:"PRODUCT_SVC_URL"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
