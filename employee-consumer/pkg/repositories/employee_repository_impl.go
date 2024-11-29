package repositories

import (
	"database/sql"
	"log"

	"github.com/MarkoLuna/EmployeeConsumer/pkg/models"
)

type EmployeeRepositoryImpl struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB, initDb bool) EmployeeRepository {
	repo := &EmployeeRepositoryImpl{db}

	if initDb {
		repo.CreateTable()
	}

	return repo
}

func (er EmployeeRepositoryImpl) CreateTable() {

	sqlStatement := `CREATE TABLE IF NOT EXISTS employees (
		id_employee TEXT PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		second_last_name TEXT,
		date_of_birth      DATE,
		date_of_employment DATE,
		status TEXT
	  );
	  `

	_, err := er.db.Exec(sqlStatement)
	if err != nil {
		log.Println("Unable to create the table:")
		log.Println(err)
	}
}

func (er EmployeeRepositoryImpl) Create(e models.Employee) (*models.Employee, error) {

	sqlStatement := `
		INSERT INTO employees (
			id_employee,
			first_name, 
			last_name, 
			second_last_name, 
			date_of_birth,
			date_of_employment, 
			status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id_employee`

	err := er.db.QueryRow(sqlStatement, e.Id, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status).Scan(&e.Id)

	if err != nil {
		return nil, err
	}
	log.Println("New record ID is:", e.Id)

	return &e, nil
}

func (er EmployeeRepositoryImpl) FindAll() ([]models.Employee, error) {

	rowsRs, err := er.db.Query(`SELECT id_employee,
							first_name,
							last_name,
							second_last_name,
							date_of_birth,
							date_of_employment,
							status FROM employees`)

	employeesSlice := make([]models.Employee, 0)

	if err != nil {
		log.Println("Error getting employees: ")
		log.Println(err)
		return employeesSlice, err
	}

	for rowsRs.Next() {
		employee := models.Employee{}
		err := rowsRs.Scan(&employee.Id, &employee.FirstName,
			&employee.LastName, &employee.SecondLastName,
			&employee.DateOfBirth, &employee.DateOfEmployment, &employee.Status)

		if err != nil {
			log.Println(err)
			continue
		}
		employeesSlice = append(employeesSlice, employee)
	}

	if err = rowsRs.Err(); err != nil {
		log.Println(err)
		return employeesSlice, err
	}

	return employeesSlice, nil
}

func (er EmployeeRepositoryImpl) ExistsById(ID string) bool {
	// TODO add impl
	return false
}

func (er EmployeeRepositoryImpl) FindById(ID string) (models.Employee, error) {

	var employee models.Employee
	userSql := `select
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
					id_employee = $1`

	err := er.db.QueryRow(userSql, ID).Scan(&employee.Id,
		&employee.FirstName, &employee.LastName,
		&employee.SecondLastName, &employee.DateOfBirth,
		&employee.DateOfEmployment, &employee.Status)

	if err != nil {
		log.Printf("Failed to execute query: %s", err)
		return employee, err
	}

	return employee, nil
}

func (er EmployeeRepositoryImpl) DeleteById(ID string) (int64, error) {

	sqlStatement := `DELETE FROM employees WHERE id_employee = $1;`

	res, err := er.db.Exec(sqlStatement, ID)

	if err != nil {
		log.Println("Unable to delete the row:")
		log.Println(err)
		return 0, err
	}

	count, _ := res.RowsAffected()
	return count, nil
}

func (er EmployeeRepositoryImpl) Update(e models.Employee) (int64, error) {

	sqlStatement := `
		UPDATE employees SET 
			first_name = $2, 
			last_name = $3, 
			second_last_name = $4,
			date_of_birth = $5,
			date_of_employment = $6,
			status = $7 
		WHERE id_employee = $1;`

	res, err := er.db.Exec(sqlStatement, e.Id, e.FirstName, e.LastName, e.SecondLastName,
		e.DateOfBirth, e.DateOfEmployment, e.Status)

	if err != nil {
		log.Println("Unable to update the row:")
		log.Println(err)
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Println("Unable to update the row:")
		log.Println(err)
		return 0, err
	}
	log.Printf("Rows affected: %d\n", count)

	return count, nil
}
