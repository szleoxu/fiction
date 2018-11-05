package config

import (
	"github.com/Unknwon/goconfig"
	"fmt"
	"os"
)

func GetConfigFile() *goconfig.ConfigFile {
	cfg, configErr := goconfig.LoadConfigFile("config/config.ini")
	if configErr != nil {
		fmt.Println(configErr)
		os.Exit(1)
	}
	return cfg
}


func GetValue(key string,value string) string {
	cfg := GetConfigFile()
	urlDDXS, _ := cfg.GetValue(key, value)
	return urlDDXS
}
