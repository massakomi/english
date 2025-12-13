package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type Word struct {
	Id          int
	Word        string
	Translate   string
	Comment     sql.NullString
	DateAdded   time.Time
	DateUpdated sql.NullTime
	DateRemind  sql.NullTime
	IdCategory  int
	DateResult  sql.NullTime
}

func GetWords(database *sqlx.DB, where string) []Word {
	if !strings.Contains(where, "order by") {
		where += ` order by id`
	}
	s := fmt.Sprintf(`select * from english where %v`, where)
	data := GetWordsBySql(database, s, func(rows *sql.Rows, p *Word) error {
		return rows.Scan(&p.Id, &p.Word, &p.Translate, &p.Comment, &p.DateAdded, &p.DateUpdated, &p.DateRemind, &p.IdCategory)
	})
	return data
}

func GetWordsWithDateResult(database *sqlx.DB) []Word {
	s := `
		select e.*, MAX(el.date_added) as date_result
    	from english e
    	left join english_log el ON el.id_word=e.id
    	group by e.id
    	order by id`
	data := GetWordsBySql(database, s, func(rows *sql.Rows, p *Word) error {
		return rows.Scan(&p.Id, &p.Word, &p.Translate, &p.Comment, &p.DateAdded, &p.DateUpdated, &p.DateRemind, &p.IdCategory, &p.DateResult)
	})
	return data
}

// GetWordsBySql Общая механика выборки из таблицы
func GetWordsBySql(database *sqlx.DB, sql string, scanner func(rows *sql.Rows, p *Word) error) []Word {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []Word
	for rows.Next() {
		p := Word{}
		err := scanner(rows, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
