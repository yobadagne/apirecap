package util

import (
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
	"github.com/spf13/viper"
	"github.com/yobadagne/user_registration/model"
)

type Config struct {
	Access_key  string `mapstructure:"ACCESS_KEY"`
	Refersh_key string `mapstructure:"REFRESH_KEY"`
}

func LoadConfig(c *gin.Context, path string) ( config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// so that it can read from environment variable
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		err = errorx.Decorate(err, "Error in viper config")
		Logger.Error("Error in viper config")
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		err = errorx.Decorate(err, "Error in viper unmarshal")
		Logger.Error("Error in viper unmarshal")
		c.Set(model.Error_type, model.INTERNAL_SERVER_ERROR)
		return
	}
	return
}
