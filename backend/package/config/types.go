package config

// Config configuration yaml structure
type Config struct {
	Delivery   `yaml:"delivery" mapstructure:",squash"`
	View       `yaml:"view" mapstructure:",squash"`
	Repository `yaml:"repository" mapstructure:",squash"`
	Users      `yaml:"users" mapstructure:",squash"`
}

// Delivery injected to backend/delivery
type Delivery struct {
	Port       int    `yaml:"http_port" mapstructure:"HTTP_PORT"`
	CookieName string `yaml:"cookie_name" mapstructure:"COOKIE_NAME"`
	ByPass     bool   `yaml:"bypass"`
}

// View injected to backend/view
type View struct {
	Port int    `yaml:"http_port"`
	Path string `yaml:"static_path" mapstructure:"STATIC_PATH"`
}

// Repository injected to backend/repository
type Repository struct {
	DatabaseHost     string `yaml:"db_host" mapstructure:"DB_HOST"`
	DatabaseUser     string `yaml:"db_user" mapstructure:"DB_USER"`
	DatabaseName     string `yaml:"db_name" mapstructure:"DB_NAME"`
	DatabasePassword string `yaml:"db_pass" mapstructure:"DB_PASS"`
	DefaultPassword  string `yaml:"default_pass" mapstructure:"DEFAULT_PASS"`
}

// Users injected to backend/usecase/users
type Users struct {
	JWTKey       string `yaml:"jwt_key" mapstructure:"JWT_KEY"`
	TokenExpiry  int64  `yaml:"token_expiry" mapstructure:"TOKEN_EXPIRY"`
	RefreshToken int64  `yaml:"refresh_token" mapstructure:"REFRESH_TOKEN"`
}
