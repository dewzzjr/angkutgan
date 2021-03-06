package database

import (
	"database/sql"
	"time"

	"github.com/dewzzjr/angkutgan/backend/model"
)

// NullInt64 use null if default value
func NullInt64(value int64) (result sql.NullInt64) {
	if value != 0 {
		result = sql.NullInt64{
			Int64: value,
			Valid: true,
		}
	}
	return
}

// NullString use null if default value
func NullString(value string) (result sql.NullString) {
	if value != "" {
		result = sql.NullString{
			String: value,
			Valid:  true,
		}
	}
	return
}

// NullInt use null if default value
func NullInt(value int) (result sql.NullInt64) {
	if value != 0 {
		result = sql.NullInt64{
			Int64: int64(value),
			Valid: true,
		}
	}
	return
}

// NullTime use null if default value
func NullTime(value string) (result sql.NullTime) {
	if date, err := time.Parse(model.DateFormat, value); err == nil {
		result = sql.NullTime{
			Time:  date,
			Valid: true,
		}
	}
	return
}
