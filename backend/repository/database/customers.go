package database

import (
	"context"
	"reflect"
	"strings"

	"github.com/dewzzjr/angkutgan/backend/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const qGetListCustomers = `SELECT
	code, 
	name, 
	type, 
	address, 
	phone, 
	COALESCE(nik, ''),
	COALESCE(role, ''),
	COALESCE(group_name, '')
FROM
	customers 
ORDER BY 
	create_time DESC
LIMIT ? OFFSET ?
`

// GetListCustomers using pagination
func (d *Database) GetListCustomers(ctx context.Context, limit, offset int) (customers []model.Customer, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetListCustomers, limit, offset); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%d, %d]", limit, offset)
		return
	}
	if customers, err = d.GetProjectsByCodes(ctx, rows); err != nil {
		err = errors.Wrapf(err, "GetProjectsByCodes [%d, %d]", limit, offset)
	}
	return
}

const qGetListCustomersByKeyword = `SELECT
	code, 
	name, 
	type, 
	address, 
	phone, 
	COALESCE(nik, ''),
	COALESCE(role, ''),
	COALESCE(group_name, '')
FROM
	customers 
WHERE
	code = ? OR UPPER(name) LIKE CONCAT('%', ?, '%')
ORDER BY 
	create_time DESC
LIMIT ? OFFSET ?
`

// GetListCustomersByKeyword by keyword using pagination
func (d *Database) GetListCustomersByKeyword(ctx context.Context, keyword string, limit, offset int) (customers []model.Customer, err error) {
	var rows *sqlx.Rows
	if rows, err = d.DB.QueryxContext(ctx, qGetListCustomersByKeyword,
		strings.ToUpper(keyword),
		strings.ToUpper(keyword),
		limit,
		offset,
	); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%s, %d, %d]", keyword, limit, offset)
		return
	}

	if customers, err = d.GetProjectsByCodes(ctx, rows); err != nil {
		err = errors.Wrapf(err, "GetProjectsByCodes [%d, %d]", limit, offset)
	}
	return
}

const qGetProjectByCodes = `SELECT
	code, 
	name, 
	location
FROM
	project 
WHERE
	code IN (?)
ORDER BY 
	code ASC
`

// GetProjectsByCodes bulk multiple code
func (d *Database) GetProjectsByCodes(ctx context.Context, rows *sqlx.Rows) (customers []model.Customer, err error) {
	var i int
	index := make(map[string]int)
	customers = make([]model.Customer, 0)
	for rows.Next() {
		var customer model.Customer
		if err = rows.StructScan(&customer); err != nil {
			err = errors.Wrapf(err, "StructScan")
			continue
		}
		customer.Projects = make([]model.Project, 0)
		customers = append(customers, customer)
		index[customer.Code] = i
		i++
	}

	if len(customers) == 0 {
		return
	}

	q, in, _ := sqlx.In(qGetProjectByCodes, reflect.ValueOf(index).MapKeys())
	if rows, err = d.DB.QueryxContext(ctx, q, in...); err != nil {
		err = errors.Wrapf(err, "QueryxContext [%v]", in)
		return
	}

	for rows.Next() {
		var project model.Project
		var code string
		if err = rows.Scan(
			&code,
			&project.Name,
			&project.Location,
		); err != nil {
			err = errors.Wrapf(err, "Scan [%v]", in)
			continue
		}
		if i, ok := index[code]; ok {
			customers[i].Projects = append(customers[i].Projects, project)
		}
	}
	return
}

const qGetCustomerDetail = `SELECT
	code, 
	name, 
	type, 
	address, 
	phone, 
	COALESCE(nik, ''),
	COALESCE(role, ''),
	COALESCE(group_name, '')
FROM
	customers 
WHERE
	code = ?
`

// GetCustomerDetail get customer detail by code
func (d *Database) GetCustomerDetail(ctx context.Context, code string) (customer model.Customer, err error) {
	if err = d.DB.QueryRowxContext(ctx, qGetCustomerDetail, code).StructScan(&customer); err != nil {
		err = errors.Wrapf(err, "QueryRowxContext [%s]", code)
		return
	}
	if customer.Projects, err = d.GetProjects(ctx, code); err != nil {
		err = errors.Wrapf(err, "GetProjects [%s]", code)
	}
	return
}

const qGetProjects = `SELECT
	id,
	name, 
	location
FROM
	project 
WHERE
	code = ?
`

// GetProjects list of projects by code customer
func (d *Database) GetProjects(ctx context.Context, code string) (projects []model.Project, err error) {
	projects = make([]model.Project, 0)
	if err = d.DB.SelectContext(ctx, &projects, qGetProjects, code); err != nil {
		err = errors.Wrapf(err, "SelectContext [%s]", code)
	}
	return
}

const qUpdateInsertCustomer = `INSERT
INTO
	customers (
		code,
		name,
		type, 
		address, 
		phone, 
		nik,
		role,
		group_name,
		modified_by
	)
VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? ) ON DUPLICATE KEY
UPDATE 
	name = ?, 
	type = ?, 
	address = ?, 
	phone = ?, 
	nik = ?, 
	role = ?, 
	group_name = ?, 
	modified_by = ?, 
	update_time = CURRENT_TIMESTAMP
`

// UpdateInsertCustomer insert customer or update if exists
func (d *Database) UpdateInsertCustomer(ctx context.Context, customer model.Customer, actionBy int64) (err error) {
	if _, err = d.DB.ExecContext(ctx, qUpdateInsertCustomer,
		// INSERT
		customer.Code,
		customer.Name,
		customer.Type,
		customer.Address,
		customer.Phone,
		NullString(customer.NIK),
		NullString(customer.Role),
		NullString(customer.GroupName),
		NullInt64(actionBy),
		// UPDATE
		customer.Name,
		customer.Type,
		customer.Address,
		customer.Phone,
		NullString(customer.NIK),
		NullString(customer.Role),
		NullString(customer.GroupName),
		NullInt64(actionBy),
	); err != nil {
		err = errors.Wrapf(err, "ExecContext [%v]", customer)
	}
	return
}

const (
	qDeleteCustomer = `DELETE
FROM
	customers
WHERE
	code = ?
`
	qDeleteProjects = `DELETE
FROM
	projects
WHERE
	code = ?
`
)

// DeleteCustomer delete customer including project
func (d *Database) DeleteCustomer(ctx context.Context, code string) (err error) {
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteCustomer, code); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s]", code)
		_ = tx.Rollback()
		return
	}
	if _, err = tx.ExecContext(ctx, qDeleteProjects, code); err != nil {
		err = errors.Wrapf(err, "ExecContext [%s]", code)
		_ = tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%s]", code)
	}
	return
}

const (
	qDeleteProject = `DELETE
FROM
	projects
WHERE
	id = ? AND code = ?
`
	qInsertProject = `INSERT
INTO
	projects (
		code,
		name,
		location,
		modified_by
	)
VALUES ( ?, ?, ?, ? )
`
)

// InsertDeleteProject insert and delete project transaction
func (d *Database) InsertDeleteProject(ctx context.Context, code string, insert []model.Project, delete []int64, actionBy int64) (err error) {
	var tx *sqlx.Tx
	if tx, err = d.DB.BeginTxx(ctx, nil); err != nil {
		err = errors.Wrap(err, "BeginTxx")
		return
	}
	for _, d := range delete {
		if _, err = tx.ExecContext(ctx, qDeleteProject, d, code); err != nil {
			err = errors.Wrapf(err, "ExecContext [%s, %d]", code, d)
			_ = tx.Rollback()
			return
		}
	}
	for _, i := range insert {
		if _, err = tx.ExecContext(ctx, qInsertProject,
			code,
			i.Name,
			i.Location,
			NullInt64(actionBy),
		); err != nil {
			err = errors.Wrapf(err, "ExecContext [%s, %v]", code, i)
			_ = tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrapf(err, "Commit [%s, %v, %v]", code, delete, insert)
	}
	return
}
