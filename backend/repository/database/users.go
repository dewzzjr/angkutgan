package database

import (
	"context"
	"database/sql"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const qVerifyUser = `SELECT
	password
FROM
	users
WHERE
	username = ?
LIMIT 1
`

// VerifyUser check username password
func (d *Database) VerifyUser(ctx context.Context, username, password string) (ok bool, err error) {
	var dbpassword string
	if err = d.DB.QueryRowxContext(ctx, qVerifyUser, username).Scan(&dbpassword); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrapf(err, "QueryRowxContext [%s]", username)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password)); err != nil {
		err = nil
		return
	}
	ok = true
	return
}

const (
	qCreateUser = `INSERT
INTO
	users (
		username,
		password,
		status
	)
VALUES ( ?, ?, 'ACTIVE' )
`
	qCreateInfo = `INSERT
INTO
	user_info (
		user_id,
		nik,
		ktp,
		fullname,
		address,
		phone,
		birthdate,
		religion,
		modified_by
	)
VALUES ( LAST_INSERT_ID(), ?, ?, ?, ?, ?, ?, ?, ? )
`
)

// CreateUser insert new user information
func (d *Database) CreateUser(ctx context.Context, data model.UserInfo, actionBy int64) (err error) {
	var password []byte
	if password, err = bcrypt.GenerateFromPassword([]byte(d.Config.DefaultPassword), 10); err != nil {
		err = errors.Wrapf(err, "GenerateFromPassword [%s]", d.Config.DefaultPassword)
		return
	}
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if _, err = tx.ExecContext(ctx, qCreateUser, data.Username, string(password)); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s]", data.Username)
		tx.Rollback()
		return
	}

	if _, err = tx.ExecContext(ctx, qCreateInfo,
		data.NIK,
		data.KTP,
		data.Fullname,
		data.Address,
		data.Phone,
		data.Birthdate,
		data.Religion,
		NullInt64(actionBy),
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%v]", data)
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%v]", data)
	}
	return
}

const qEditInfo = `UPDATE user_info
SET
	nik = ?,
	ktp = ?,
	fullname = ?,
	address = ?,
	phone = ?,
	birthdate = ?,
	religion = ?,
	modified_by = ?,
	update_time = CURRENT_TIMESTAMP
WHERE
	id = ? OR user_id = ?
`

// EditUser update user information
func (d *Database) EditUser(ctx context.Context, data model.UserInfo, actionBy int64) (err error) {
	if _, err = d.DB.ExecContext(ctx, qEditInfo,
		data.NIK,
		data.KTP,
		data.Fullname,
		data.Address,
		data.Phone,
		data.Birthdate,
		data.Religion,
		actionBy,
		data.InfoID,
		data.ID,
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%v]", data)
	}
	return
}

const qChangePassword = `UPDATE
	users
SET
	password = ?
WHERE
	username = ?
`

// ChangePassword change user password
func (d *Database) ChangePassword(ctx context.Context, username, newPassword string) (err error) {
	var password []byte
	if password, err = bcrypt.GenerateFromPassword([]byte(newPassword), 10); err != nil {
		err = errors.Wrapf(err, "GenerateFromPassword [%s]", newPassword)
		return
	}
	if _, err = d.DB.ExecContext(ctx, qEditInfo, string(password), username); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s, %s]", username, password)
	}
	return
}

const qGetUserLogin = `SELECT
	user_id, fullname
FROM
	users
JOIN
	user_info ON user_info.user_id = users.id AND users.username = ?
`

// GetUserLogin get user login information
func (d *Database) GetUserLogin(ctx context.Context, claim *model.Claims) (err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetUserLogin, claim.Username).Scan(&claim.UserID, &claim.Fullname); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%s]", claim.Username)
		return
	}
	if claim.UAM, err = d.GetUAM(ctx, claim.UserID); err != nil {
		err = errors.Wrap(err, "GetUAM")
	}
	return
}

const qGetUAM = `SELECT
	code, user_id
FROM
	users
JOIN
	user_access ON users.id = user_access.user_id AND users.id = ?
JOIN
	access ON user_access.access_id = access.id
`

// GetUAM get user access by username
func (d *Database) GetUAM(ctx context.Context, uid int64) (access []string, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetUAM, uid); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%d]", uid)
		return
	}
	for rows.Next() {
		var uam string
		if err = rows.Scan(&uam, &uid); err != nil {
			err = errors.Wrapf(err, "Scan [%d]", uid)
			continue
		}
		access = append(access, uam)
	}
	return
}

const qIsValidUsername = `SELECT username
FROM
	users
WHERE
	username = ?
`

// IsValidUsername check is username is valid or not to be use
func (d *Database) IsValidUsername(ctx context.Context, username string) (bool, error) {
	err := d.DB.QueryRowxContext(ctx, qVerifyUser, username).Scan(&username)
	if err == nil {
		return false, nil
	}
	if err != sql.ErrNoRows {
		err = errors.Wrapf(err, "QueryRowxContext [%s]", username)
		return false, err
	}
	return true, nil
}
