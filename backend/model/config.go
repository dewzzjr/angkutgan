package model

// Delivery injected to backend/delivery
type Delivery struct {
	Port       int    `yaml:"http_port"`
	StaticPath string `yaml:"static_path"`
	CookieName string `yaml:"cookie_name"`
}

// Repository injected to backend/repository
type Repository struct {
	DatabaseHost     string `yaml:"db_host"`
	DatabaseUser     string `yaml:"db_user"`
	DatabaseName     string `yaml:"db_name"`
	DatabasePassword string `yaml:"db_pass"`
}

// Users injected to backend/usecase/users
type Users struct {
	JWTKey       string `yaml:"jwt_key"`
	TokenExpiry  int64  `yaml:"token_expiry"`
	RefreshToken int64  `yaml:"refresh_token"`
}
