// Copyright © 2020. All rights reserved.
package repositories

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/constants"
	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"

	"github.com/stretchr/testify/assert"
)

var e = &models.Employee{
	Id:               "1",
	FirstName:        "Marcos",
	LastName:         "Luna",
	SecondLastName:   "Valdez",
	DateOfBirth:      time.Date(1994, time.April, 25, 8, 0, 0, 0, time.UTC),
	DateOfEmployment: time.Now().UTC(),
	Status:           constants.ACTIVE,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestEmployeeRepositoryImpl_FindById(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `select
				id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			from
				employees
			where
				id_employee \\?`

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"}).
		AddRow(e.Id, e.FirstName, e.LastName, e.SecondLastName, e.DateOfBirth, e.DateOfEmployment, e.Status)

	mock.ExpectQuery(query).WithArgs(e.Id).WillReturnRows(rows)

	employee, err := repo.FindById(e.Id)
	assert.NotNil(t, employee)
	assert.NoError(t, err)
}

func TestEmployeeRepositoryImpl_FindByIdError(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `select
				id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status
			from
				employees
			where
				id_employee `

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"})

	mock.ExpectQuery(query).WithArgs(e.Id).WillReturnRows(rows)

	employee, err := repo.FindById(e.Id)
	assert.Empty(t, employee)
	assert.Error(t, err)

	err1 := errors.New("sql: no rows in result set")
	assert.Equal(t, err1, err)
}

func TestEmployeeRepositoryImpl_FindAll(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `SELECT id_employee,
				first_name,
				last_name,
				second_last_name,
				date_of_birth,
				date_of_employment,
				status 
			  FROM employees`

	rows := sqlmock.NewRows([]string{"id_employee", "first_name", "last_name", "second_last_name",
		"date_of_birth", "date_of_employment", "status"}).
		AddRow(e.Id, e.FirstName, e.LastName, e.SecondLastName, e.DateOfBirth, e.DateOfEmployment, e.Status)

	mock.ExpectQuery(query).WillReturnRows(rows)

	employees, err := repo.FindAll()

	assert.NotEmpty(t, employees)
	assert.NoError(t, err)
	assert.Len(t, employees, 1)
}

func TestEmployeeRepositoryImpl_Create(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `
	INSERT INTO employees \(
		id_employee\,
		first_name\,
		last_name\,
		second_last_name\,
		date_of_birth\,
		date_of_employment\,
		status\)
	VALUES \(\$1\, \$2\, \$3\, \$4\, \$5\, \$6\, \$7\)
	RETURNING id_employee`

	rows := sqlmock.NewRows([]string{"id_employee"}).AddRow(1)

	mock.ExpectQuery(query).WithArgs(e.Id, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status).WillReturnRows(rows)

	_, err := repo.Create(*e)
	assert.NoError(t, err)
}

func TestEmployeeRepositoryImpl_Update(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `
		UPDATE employees SET
			first_name = \$2\,
			last_name = \$3\,
			second_last_name = \$4\,
			date_of_birth = \$5\,
			date_of_employment = \$6\,
			status = \$7
		WHERE id_employee = \$1;
	`

	mock.ExpectExec(query).WithArgs(e.Id, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status).WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := repo.Update(*e)
	assert.NoError(t, err)
}

func TestEmployeeRepositoryImpl_UpdateErr(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `
		UPDATE employees SET
			first_name = \$2\,
			last_name = \$3\,
			second_last_name = \$4\,
			date_of_birth = \$5\,
			date_of_employment = \$6\,
			status = \$7
		WHERE id_employee = \$1;
	`

	mock.ExpectExec(query).WithArgs(e.Id, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status).WillReturnResult(sqlmock.NewResult(0, 0))

	count, err := repo.Update(*e)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(0))
}

func TestEmployeeRepositoryImpl_DeleteById(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `DELETE FROM employees WHERE id_employee = \$1\;`
	mock.ExpectExec(query).WithArgs(e.Id).WillReturnResult(sqlmock.NewResult(0, 1))

	count, err := repo.DeleteById(e.Id)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(1))
}

func TestEmployeeRepositoryStub_DeleteByIdError(t *testing.T) {
	db, mock := NewMock()
	repo := &EmployeeRepositoryImpl{db}
	defer func() {
		db.Close()
	}()

	query := `DELETE FROM employees WHERE id_employee = \$1\;`
	mock.ExpectExec(query).WithArgs(e.Id).WillReturnResult(sqlmock.NewResult(0, 0))

	count, err := repo.DeleteById(e.Id)
	assert.NoError(t, err)
	assert.Equal(t, count, int64(0))
}
