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
	Port       int    `yaml:"http_port" ignored:"true"`
	CookieName string `yaml:"cookie_name" envconfig:"COOKIE_NAME"`
	ByPass     bool   `yaml:"bypass"`
}

// View injected to backend/view
type View struct {
	ViewPort int    `yaml:"-" ignored:"true"`
	Path     string `yaml:"static_path" envconfig:"STATIC_PATH"`
}

// Repository injected to backend/repository
type Repository struct {
	DatabaseEnvKey   string `yaml:"-" envconfig:"DB_KEY"`
	DatabaseURL      string `yaml:"db_url" ignored:"true"`
	DatabaseHost     string `yaml:"db_host" envconfig:"DB_HOST"`
	DatabaseUser     string `yaml:"db_user" envconfig:"DB_USER"`
	DatabaseName     string `yaml:"db_name" envconfig:"DB_NAME"`
	DatabasePassword string `yaml:"db_pass" envconfig:"DB_PASS"`
	DefaultPassword  string `yaml:"default_pass" envconfig:"DEFAULT_PASS"`
}

// Users injected to backend/usecase/users
type Users struct {
	JWTKey       string `yaml:"jwt_key" envconfig:"JWT_KEY"`
	TokenExpiry  int64  `yaml:"token_expiry" envconfig:"TOKEN_EXPIRY"`
	RefreshToken int64  `yaml:"refresh_token" envconfig:"REFRESH_TOKEN"`
}
