package main

import (
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", `dbname=root host=localhost user=postgres password=password`)
	if err != nil {
		panic(err)
	}

	boil.SetDB(db)

	RunOrchestrator()
}
