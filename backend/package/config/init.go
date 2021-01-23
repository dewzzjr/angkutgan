package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var config *Config

func init() {
	if config == nil {
		env := os.Getenv("ENVIRONMENT")
		switch env {
		case "PRODUCTION", "STAGING":
			loadEnv(env)
		default:
			loadFile()
		}
		if config == nil {
			log.Fatalln("no config loaded")
		}
		validate()
	}
}

func loadFile() {
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
	}
}

func loadEnv(env string) {
	const Prefix = "APP"
	config = &Config{}
	if err := envconfig.Process(Prefix, &config.Delivery); err != nil {
		log.Fatalln("config.Delivery", env, err)
	}
	if err := envconfig.Process(Prefix, &config.Users); err != nil {
		log.Fatalln("config.Users", env, err)
	}
	if err := envconfig.Process(Prefix, &config.View); err != nil {
		log.Fatalln("config.View", env, err)
	}
	if err := envconfig.Process(Prefix, &config.Repository); err != nil {
		log.Fatalln("config.Repository", env, err)
	}
	if port, ok := os.LookupEnv("PORT"); ok {
		config.Port, _ = strconv.Atoi(port)
	}
	config.Repository.DatabaseURL = os.Getenv(config.DatabaseEnvKey)
}

func validate() {
	if config.ViewPort == 0 {
		config.ViewPort = config.Port
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

// Get configuration
func Get() *Config {
	if config == nil {
		log.Fatal("no config found")
	}
	return config
}
