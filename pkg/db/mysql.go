package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"maps"
	"slices"
	"strconv"
	"time"
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

// GetFirstVal получаем первое значение первой строки в виде string
func GetFirstVal(sql string, database *sqlx.DB) string {
	row := GetData(sql, database)
	if len(row) == 0 {
		return ""
	}
	key := slices.Sorted(maps.Keys(row[0]))[0]
	if row[0][key] == nil {
		return ""
	}

	var dt string
	switch v := row[0]["dt"].(type) {
	case int:
		dt = strconv.FormatInt(int64(row[0]["dt"].(int)), 10)
	case string:
		dt = row[0]["dt"].(string)
	case time.Time:
		dt = row[0]["dt"].(time.Time).Format("2006-01-02 15:04:05")
	default:
		log.Panicf("I don't know about type %T!\n", v)
	}
	return dt
}
