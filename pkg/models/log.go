package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Log struct {
	Id        int
	IdWord    int
	DateAdded time.Time
	Result    int
	Time      float64
}

func GetLogs(database *sqlx.DB) []Log {
	s := `select * from english_log order by id desc`
	data := GetLogsBySql(database, s, func(rows *sql.Rows, p *Log) error {
		return rows.Scan(&p.Id, &p.IdWord, &p.DateAdded, &p.Result, &p.Time)
	})
	return data
}

// GetLogsBySql Общая механика выборки из таблицы
func GetLogsBySql(database *sqlx.DB, sql string, scanner func(rows *sql.Rows, p *Log) error) []Log {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []Log
	for rows.Next() {
		p := Log{}
		err := scanner(rows, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
