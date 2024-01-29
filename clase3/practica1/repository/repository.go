package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrStorageProductInternal  = errors.New("internal storage product error")
	ErrStorageProductNotFound  = errors.New("storage product not found")
	ErrStorageProductNotUnique = errors.New("storage product not unique")
)

type Product struct {
	Id         int
	Name       string
	Quantity   int
	CodeValue  string
	Published  bool
	Expiration time.Time
	Price      float64
}

type ProductSQL struct {
	Id         sql.NullInt32
	Name       sql.NullString
	Quantity   sql.NullInt32
	CodeValue  sql.NullString
	Published  sql.NullString
	Expiration sql.NullTime
	Price      sql.NullFloat64
}

// NewImplStorageProductMySQL returns new ImplStorageProductMySQL
func NewImplStorageProductMySQL(db *sql.DB) *ImplStorageProductMySQL {
	return &ImplStorageProductMySQL{db: db}
}

// ImplStorageProductMySQL is an implementation of StorageProduct interface
type ImplStorageProductMySQL struct {
	db *sql.DB
}

// GetOne returns one product by id
func (impl *ImplStorageProductMySQL) GetOne(id int) (p *Product, err error) {
	// query
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	row := stmt.QueryRow(id)
	if row.Err() != nil {
		err = row.Err()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = fmt.Errorf("%w. %v", ErrStorageProductNotFound, row.Err())
		default:
			err = fmt.Errorf("%w. %v", ErrStorageProductInternal, row.Err())
		}

		return
	}

	// scan row
	var product ProductSQL
	err = row.Scan(&product.Id, &product.Name, &product.Quantity, &product.CodeValue, &product.Published, &product.Expiration, &product.Price)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// serialization
	p = new(Product)
	p.Name = product.Name.String
	p.Id = int(product.Id.Int32)
	p.Quantity = int(product.Quantity.Int32)
	p.CodeValue = product.CodeValue.String
	p.Published = product.Published.Valid
	p.Expiration = product.Expiration.Time
	p.Price = product.Price.Float64

	return
}

// Store stores product
func (impl *ImplStorageProductMySQL) Store(p *Product) (err error) {
	// deserialize
	var product ProductSQL

	product.Name.Valid = true
	product.Name.String = (*p).Name

	product.Quantity.Valid = true
	product.Quantity.Int32 = int32((*p).Quantity)

	product.CodeValue.Valid = true
	product.CodeValue.String = (*p).CodeValue

	product.Published.Valid = true
	product.Published.String = strconv.FormatBool((*p).Published)

	product.Expiration.Valid = true
	product.Expiration.Time = (*p).Expiration

	product.Price.Valid = true
	product.Price.Float64 = (*p).Price

	// query
	query := "INSERT INTO products (name, quantity, code_value, is_published, expiration, price) VALUES (?, ?, ?, ?, ?, ?)"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(product.Name, product.Quantity, product.CodeValue, product.Published, product.Expiration, product.Price)
	if err != nil {
		errMySQL, ok := err.(*mysql.MySQLError)
		if ok {
			switch errMySQL.Number {
			case 1062:
				err = fmt.Errorf("%w. %v", ErrStorageProductNotUnique, err)
			default:
				err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
			}

			return
		}

		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	// get last insert id
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	(*p).Id = int(lastInsertID)

	return
}

// Update updates product
func (impl *ImplStorageProductMySQL) Update(p *Product) (err error) {
	// deserialize
	var product ProductSQL
	product.Name.Valid = true
	product.Name.String = (*p).Name

	product.Quantity.Valid = true
	product.Quantity.Int32 = int32((*p).Quantity)

	product.CodeValue.Valid = true
	product.CodeValue.String = (*p).CodeValue

	product.Published.Valid = true
	product.Published.String = strconv.FormatBool((*p).Published)

	product.Expiration.Valid = true
	product.Expiration.Time = (*p).Expiration

	product.Price.Valid = true
	product.Price.Float64 = (*p).Price

	// query
	query := "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(product.Name, product.Quantity, product.CodeValue, product.Published, product.Expiration, product.Price, p.Id)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	return
}

// Delete deletes product by id
func (impl *ImplStorageProductMySQL) Delete(id int) (err error) {
	// query
	query := "DELETE FROM products WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(id)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	return
}
