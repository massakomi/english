package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

func Connect() *sqlx.DB {
	connStr := "host=localhost user=tester password=tester port=5431 dbname=tester sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	//if err = db.Ping(); err != nil {
	//	panic(err)
	//}
	return db
}

func Fields(db *sqlx.DB, table string) []string {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %v", table))
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	defer rows.Close()
	return columns
}

func GetData(sql string, db *sqlx.DB) []map[string]interface{} {
	rows, err := db.Queryx(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		rowMap := make(map[string]interface{})
		if err := rows.MapScan(rowMap); err != nil {
			panic(err)
		}
		results = append(results, rowMap)
	}
	return results
}

// Singleton represents the type for which we want a single instance.
type Singleton struct {
	data string
}

var (
	instance *Singleton
	once     sync.Once
)

// GetInstance returns the single instance of the Singleton.
func GetInstance() *Singleton {
	once.Do(func() {
		// This code block will be executed only once, even with concurrent calls.
		instance = &Singleton{data: "Initialized Data"}
		fmt.Println("Singleton instance created.")
	})
	return instance
}

// ExampleMethod demonstrates a method on the Singleton instance.
func (s *Singleton) ExampleMethod() {
	fmt.Printf("Singleton data: %s\n", s.data)
}
