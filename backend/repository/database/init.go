package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dewzzjr/angkutgan/backend/package/config"
	"github.com/jmoiron/sqlx"
)

// Database repository object
type Database struct {
	DB     iDatabase
	Config config.Repository
}

// New initiate repository/database
func New(cfg config.Repository) *Database {
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabaseName)
	} else {
		cfg.DatabaseURL = fmt.Sprintf("%s%s", cfg.DatabaseURL, "?parseTime=true")
	}
	return &Database{
		DB:     sqlx.MustConnect("mysql", cfg.DatabaseURL),
		Config: cfg,
	}
}

type iDatabase interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Begin() (*sql.Tx, error)
	Beginx() (*sqlx.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Rebind(sql string) string
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}
