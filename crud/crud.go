package crud

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Employee struct { // Arbitrary non-encrypted information
	Username string
	Password string
	Position string
	Salary   uint64
}

// Should note that this example is specifically met for the above struct, otherwise issues would occur.

func sqlExec(dbpool *pgxpool.Pool, sqlString string) (pgconn.CommandTag, error) {
	result, err := dbpool.Exec(context.Background(), sqlString)
	return result, err
}

func sqlQueryRow(dbpool *pgxpool.Pool, sqlString string) (Employee, error) {
	var employee Employee
	if err := dbpool.QueryRow(context.Background(), sqlString).
		Scan(nil, &employee.Username, &employee.Password, &employee.Position, &employee.Salary); err != nil {
		return employee, err
	}
	return employee, nil
}

func SQLCreateTable(dbpool *pgxpool.Pool) (string, error) {
	var sqlString string = "DROP SEQUENCE IF EXISTS seq CASCADE;"
	_, err := sqlExec(dbpool, sqlString)
	if err != nil {
		fmt.Println("WARNING: Drop sequence failed (maybe non-existant?)")
	}

	sqlString = "DROP TABLE Employees CASCADE;"
	_, err = sqlExec(dbpool, sqlString)
	if err != nil {
		fmt.Println("WARNING: Drop table cascade for employees failed (maybe non-existant?)")
	}

	sqlString = "CREATE SEQUENCE seq INCREMENT BY 1;"
	_, err = sqlExec(dbpool, sqlString)
	if err != nil {
		return "FATAL: Sequence creation failed\n", err
	}

	sqlString = (`
		CREATE TABLE Employees(
			ID INTEGER NOT NULL DEFAULT nextval('seq'),
			Username VARCHAR(255) NOT NULL UNIQUE,
			Password VARCHAR(255) NOT NULL,
			Position VARCHAR(255),
			Salary INTEGER,
			PRIMARY KEY (ID)
		);
	`)
	_, err = sqlExec(dbpool, sqlString)
	if err != nil {
		return "FATAL: Table creation failed\n", err
	}
	return "INFO: Table creation ready\n", nil
}

func SQLCreate(dbpool *pgxpool.Pool, employee Employee) (pgconn.CommandTag, error) {
	var sqlString string = fmt.Sprintf(
		"INSERT INTO Employees(Username, Password, Position, Salary) VALUES('%s', '%s', '%s', %d);",
		employee.Username, employee.Password, employee.Position, employee.Salary)
	result, err := sqlExec(dbpool, sqlString)
	return result, err
}

func SQLRead(dbpool *pgxpool.Pool, employee Employee) (Employee, error) {
	var sqlString string = fmt.Sprintf(
		`SELECT * FROM Employees WHERE 
		Username='%v' AND Password='%v' AND Position='%v' AND SALARY=%v;`,
		employee.Username, employee.Password, employee.Position, employee.Salary)
	result, err := sqlQueryRow(dbpool, sqlString)
	return result, err
}

func SQLUpdate(dbpool *pgxpool.Pool, employee Employee, sa []string) (pgconn.CommandTag, error) {
	var sqlString string = fmt.Sprintf(
		"UPDATE Employees SET %v = %v WHERE Username = '%s';",
		sa[0], sa[1], employee.Username)
	result, err := sqlExec(dbpool, sqlString)
	return result, err
}

func SQLDelete(dbpool *pgxpool.Pool, employee Employee) (pgconn.CommandTag, error) {
	var sqlString string = fmt.Sprintf(
		"DELETE FROM Employees WHERE Username = '%s';", employee.Username)
	result, err := sqlExec(dbpool, sqlString)
	return result, err
}
