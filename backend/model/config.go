package model

// Delivery injected to backend/delivery
type Delivery struct {
	Port int `yaml:"http_port"`
}

// Repository injected to backend/repository
type Repository struct {
	DatabaseHost     string `yaml:"db_host"`
	DatabaseUser     string `yaml:"db_user"`
	DatabaseName     string `yaml:"db_name"`
	DatabasePassword string `yaml:"db_pass"`
}
