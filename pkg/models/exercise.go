package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"html/template"
	"time"
)

type Exercise struct {
	Id           int
	Page         int
	DateAdded    time.Time
	DateFinished time.Time
	DaysAgo      int
	Time         float64
	Style        template.CSS
	Errors       int
}

func GetExercises(database *sqlx.DB, where string, order string) []Exercise {
	s := `select * from english_exercise`
	if where != "" {
		s += " where " + where
	}
	s += order
	data := GetExercisesBySql(database, s, func(rows *sql.Rows, p *Exercise) error {
		return rows.Scan(&p.Id, &p.Page, &p.DateAdded, &p.DateFinished)
	})
	return data
}

// GetEnglishBookReadBySql Общая механика выборки из таблицы
func GetExercisesBySql(database *sqlx.DB, sql string, scanner func(rows *sql.Rows, p *Exercise) error) []Exercise {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []Exercise
	for rows.Next() {
		p := Exercise{}
		err := scanner(rows, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
