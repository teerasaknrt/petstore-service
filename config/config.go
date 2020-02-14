package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Debug    bool
	Env      string
	GRPCPort string
	HTTPPort string
	Firestore
}

type Firestore struct {
	URL string
}

var config Config

func InitViper(path string, env string) {
	switch env {
	case "local":
		viper.SetConfigName("local-config")
		break
	case "develop":
		viper.SetConfigName("dev-config")
		break
	case "staging":
		viper.SetConfigName("staging-config")
		break
	default:
		viper.SetConfigName("config")
		break
	}
	log.Println("running on environment :", env)
	viper.AddConfigPath(path)
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
}

func GetViper() *Config {
	return &config
}
