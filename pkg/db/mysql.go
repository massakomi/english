package db

import (
	"database/sql"
	"fmt"
)

func Connect() *sql.DB {
	connStr := "host=localhost user=tester password=tester port=5431 dbname=tester sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	//if err = db.Ping(); err != nil {
	//	panic(err)
	//}
	return db
}

func Fields(db *sql.DB, table string) []string {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %v", table))
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	defer rows.Close()
	return columns
}
