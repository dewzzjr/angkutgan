package config

import "github.com/dewzzjr/angkutgan/backend/model"

// Config configuration yaml structure
type Config struct {
	model.Delivery   `yaml:"delivery"`
	model.View       `yaml:"view"`
	model.Repository `yaml:"repository"`
	model.Users      `yaml:"users"`
}
