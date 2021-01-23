package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var config *Config

func init() {
	if config == nil {
		loadEnv()
	}
}

func load() {
	path := []string{
		"./config.yaml",
		"../config.yaml",
		"../../config.yaml",
		"./configs/angkutgan.yaml",
		"../configs/angkutgan.yaml",
		"../../configs/angkutgan.yaml",
		"/etc/conf/angkutgan.yaml",
	}
	for _, p := range path {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		if err = yaml.NewDecoder(f).Decode(&config); err != nil {
			continue
		}
		validate()
	}
	if config == nil {
		log.Println("no config loaded")
	}
}

func loadEnv() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs/")
	viper.AddConfigPath("../configs/")
	viper.SetConfigName("angkutgan")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}
	if config == nil {
		log.Println("no config loaded")
		return
	}
	validate()
}

func validate() {
	if config.View.Port == 0 {
		config.View.Port = config.Delivery.Port
	}
	config.View.Path = check(config.View.Path, 0)
}

func check(path string, tries int) string {
	if _, err := os.Stat(path); os.IsNotExist(err) && tries < 3 {
		path = "../" + strings.TrimPrefix(path, "/")
		return check(path, tries+1)
	}
	return path
}

// Load config from specific
func Load(path string) (err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return
	}
	return
}

// Get configuration
func Get() *Config {
	if config == nil {
		log.Fatal("no config found")
	}
	return config
}
