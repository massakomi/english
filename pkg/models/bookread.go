package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type EnglishBookread struct {
	id           int
	Name         string
	idBook       int
	Page         int
	DateAdded    time.Time
	DateFinished time.Time
	Current      bool
	Offset       float64
}

// Метод с встроенным запросом для выборки
func GetEnglishBookRead(database *sqlx.DB, withExercises bool, exercisesMaxDays int) []EnglishBookread {
	limit := 100
	add := ""
	if withExercises {
		add = fmt.Sprintf(`UNION (
			select id, 'упражнения' as name, 0 AS id_book, page, date_added, date_finished 
			from english_exercise 
			where date_added > date_add(NOW(), -'%v day'::interval))`, exercisesMaxDays)
	}
	sql := fmt.Sprintf(`select * from english_bookread %v order by date_added desc limit %v`, add, limit)
	dataEnglishAll := GetEnglishBookReadBySql(database, sql)
	return dataEnglishAll
}

// Общая механика выборки из таблиц книг в структуру (альтернативно выборке в словарь)
func GetEnglishBookReadBySql(database *sqlx.DB, sql string) []EnglishBookread {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []EnglishBookread
	for rows.Next() {
		p := EnglishBookread{}
		err := rows.Scan(&p.id, &p.Name, &p.idBook, &p.Page, &p.DateAdded, &p.DateFinished)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
