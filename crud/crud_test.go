package crud

import (
	"context"
	"log"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// All Test are given in sequential order in series. They must run before the other.
// Does not cover all test.

func setup() *pgxpool.Pool { // Should be noted that if the database does not exist then there is no reason to continue.
	if err := godotenv.Load(".env"); err != nil { // Loading the environment file (".env")
		log.Fatal("Failed to retrieve environment file.")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL")) // Creating a new connection pool
	if err != nil {
		dbpool.Close()
		log.Fatal("Failed to create the connection pool")
	}

	if err = dbpool.Ping(context.Background()); err != nil { // Basic ping test to ensure connectivity
		dbpool.Close()
		log.Fatal("Failed to ping the connection pool.")
	}

	return dbpool
}

func teardown() {
	// Currently has no use, but may come up later. ...Maybe.
}

func TestSQLCreateTable(t *testing.T) {
	dbpool := setup()
	defer dbpool.Close()

	got, err := SQLCreateTable(dbpool)
	if err != nil {
		t.Errorf("%v", err)
	} else if utf8.RuneCountInString(got) <= 0 {
		t.Errorf("%v want string with length > 0", got)
	}
}

func TestSQLCreate(t *testing.T) {
	dbpool := setup()
	defer dbpool.Close()

	employee := Employee{Username: "Goodbye", Password: "Hello", Position: "Give me", Salary: 12000}
	got, err := SQLCreate(dbpool, employee)
	if err != nil {
		t.Errorf("%v", err)
	} else if got.String() != "INSERT 0 1" {
		t.Errorf("%s want 'INSERT 0 1'", got.String())
	}

	teardown()
}

func TestSQLRead(t *testing.T) {
	dbpool := setup()
	defer dbpool.Close()

	employee := Employee{Username: "Goodbye", Password: "Hello", Position: "Give me", Salary: 12000}

	got, err := SQLRead(dbpool, employee)
	if err != nil {
		t.Errorf("%v", err)
	} else if got != employee {
		t.Errorf("%v want %v", got, employee)
	}

	teardown()
}

func TestSQLUpdate(t *testing.T) {
	dbpool := setup()
	defer dbpool.Close()

	employee := Employee{Username: "Goodbye", Password: "Hello", Position: "Give me", Salary: 12000}
	newPassword := "'Walnut'"

	sa := make([]string, 2)
	sa[0] = "Password"
	sa[1] = newPassword

	got, err := SQLUpdate(dbpool, employee, sa)
	if err != nil {
		t.Errorf("%v", err)
	} else if got.String() != "UPDATE 1" {
		t.Errorf("%s want UPDATE 1", got.String())
	}

	teardown()
}

func TestSQLDelete(t *testing.T) {
	dbpool := setup()
	defer dbpool.Close()

	employee := Employee{Username: "Goodbye", Password: "Hello", Position: "Give me", Salary: 12000}

	got, err := SQLDelete(dbpool, employee)
	if err != nil {
		t.Errorf("%v", err)
	} else if got.String() != "DELETE 1" {
		t.Errorf("%s want DELETE 1", got.String())
	}
}
