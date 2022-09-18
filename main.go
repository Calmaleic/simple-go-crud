package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Serenade419/basic-crud/crud"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil { // Loading the environment file (".env")
		log.Fatal(err)
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL")) // Creating a new connection pool
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	if err = dbpool.Ping(context.Background()); err != nil { // Basic ping test to ensure connectivity
		log.Fatal(err)
	}

	result, err := crud.SQLCreateTable(dbpool)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	employees := make([]crud.Employee, 3)
	employees[0] = crud.Employee{Username: "Chun-Li", Password: "FJNEIOW421934%!@$", Position: "System Administrator", Salary: 50000}
	employees[1] = crud.Employee{Username: "M. Bison", Password: "fewako419042190fs%$@!", Position: "Project Manager", Salary: 50000}
	employees[2] = crud.Employee{Username: "Guile", Password: "fewako419042190fs%$@!", Position: "Lead Developer", Salary: 50000}

	for _, employee := range employees {
		tag, err := crud.SQLCreate(dbpool, employee)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tag)
	}

	nemployee, err := crud.SQLRead(dbpool, employees[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nemployee)
}
