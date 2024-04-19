package util

import (
	"github.com/spf13/viper"
	"github.com/yobadagne/user_registration/model"
	"go.uber.org/zap"
)

type Config struct {
	Access_key  string `mapstructure:"ACCESS_KEY"`
	Refersh_key string `mapstructure:"REFRESH_KEY"`
}

func LoadConfig(path string) ( config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// so that it can read from environment variable
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		Logger.Error("Error in viper config,error while excuting util.LoadConfig()", zap.Error(err))
		err = model.ErrInternalServerErr.New("Error in viper config")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		Logger.Error("Error in viper unmarshal, error while excuting util.LoadConfig()",zap.Error(err))
		err = model.ErrInternalServerErr.New("Error in viper unmarshal")
		model.Error_type = model.INTERNAL_SERVER_ERROR
		return
	}
	return
}
