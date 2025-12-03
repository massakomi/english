package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"html/template"
	"time"
)

type Exercise struct {
	id           int
	Page         int
	DateAdded    time.Time
	DateFinished time.Time
	DaysAgo      int
	Time         float64
	Style        template.CSS
	Errors       int
}

func GetExercises(database *sqlx.DB) []Exercise {
	data := GetExercisesBySql(database, `select * from english_exercise order by date_added desc`)
	return data
}

// GetEnglishBookReadBySql Общая механика выборки из таблицы
func GetExercisesBySql(database *sqlx.DB, sql string) []Exercise {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	return scanExercises(rows)
}

func scanExercises(rows *sql.Rows) []Exercise {
	var items []Exercise
	for rows.Next() {
		p := Exercise{}
		err := rows.Scan(&p.id, &p.Page, &p.DateAdded, &p.DateFinished)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
