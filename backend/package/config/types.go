package config

// Config configuration yaml structure
type Config struct {
	Delivery   `yaml:"delivery"`
	View       `yaml:"view"`
	Repository `yaml:"repository"`
	Users      `yaml:"users"`
}

// Delivery injected to backend/delivery
type Delivery struct {
	Port       int    `yaml:"http_port"`
	CookieName string `yaml:"cookie_name"`
	ByPass     bool   `yaml:"bypass"`
}

// View injected to backend/view
type View struct {
	Port int    `yaml:"http_port"`
	Path string `yaml:"static_path"`
}

// Repository injected to backend/repository
type Repository struct {
	DatabaseHost     string `yaml:"db_host"`
	DatabaseUser     string `yaml:"db_user"`
	DatabaseName     string `yaml:"db_name"`
	DatabasePassword string `yaml:"db_pass"`
	DefaultPassword  string `yaml:"default_pass"`
}

// Users injected to backend/usecase/users
type Users struct {
	JWTKey       string `yaml:"jwt_key"`
	TokenExpiry  int64  `yaml:"token_expiry"`
	RefreshToken int64  `yaml:"refresh_token"`
}
