package config

// Config configuration yaml structure
type Config struct {
	Delivery   `yaml:",inline"`
	View       `yaml:",inline"`
	Repository `yaml:",inline"`
	Users      `yaml:",inline"`
}

// Delivery injected to backend/delivery
type Delivery struct {
	Port       int    `yaml:"http_port" envconfig:"PORT"`
	CookieName string `yaml:"cookie_name" envconfig:"APP_COOKIE_NAME"`
	ByPass     bool   `yaml:"bypass" envconfig:"-"`
}

// View injected to backend/view
type View struct {
	ViewPort int    `yaml:"-" envconfig:"PORT"`
	Path     string `yaml:"static_path" envconfig:"APP_STATIC_PATH"`
}

// Repository injected to backend/repository
type Repository struct {
	DatabaseURL      string `yaml:"db_url" envconfig:"-"`
	DatabaseHost     string `yaml:"db_host" envconfig:"APP_DB_HOST"`
	DatabaseUser     string `yaml:"db_user" envconfig:"APP_DB_USER"`
	DatabaseName     string `yaml:"db_name" envconfig:"APP_DB_NAME"`
	DatabasePassword string `yaml:"db_pass" envconfig:"APP_DB_PASS"`
	DefaultPassword  string `yaml:"default_pass" envconfig:"APP_DEFAULT_PASS"`
}

// Users injected to backend/usecase/users
type Users struct {
	JWTKey       string `yaml:"jwt_key" envconfig:"APP_JWT_KEY"`
	TokenExpiry  int64  `yaml:"token_expiry" envconfig:"APP_TOKEN_EXPIRY"`
	RefreshToken int64  `yaml:"refresh_token" envconfig:"APP_REFRESH_TOKEN"`
}
