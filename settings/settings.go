package settings

import (
	"github.com/spf13/viper"
	"log"
)

func Config() *viper.Viper {

	v := viper.New()

	v.SetConfigFile("settings/settings.yml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	return v
}
