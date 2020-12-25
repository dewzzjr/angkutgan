package database

import "database/sql"

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
