package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type EnglishWord struct {
	Id               int
	English          string
	Russian          string
	Comment          sql.NullString
	Book             string
	IdBook           int
	List             int
	Page             sql.NullInt64
	DateAdded        time.Time
	EnglishShort     sql.NullString
	EnglishShortAuto sql.NullString
}

func GetEnglishWords(database *sqlx.DB, where string) []EnglishWord {
	if !strings.Contains(where, "order by") {
		where += ` order by id`
	}
	s := fmt.Sprintf(`select * from english_words where %v`, where)
	fmt.Println(s)
	data := GetEnglishWordsBySql(database, s, func(rows *sql.Rows, p *EnglishWord) error {
		return rows.Scan(&p.Id, &p.English, &p.Russian, &p.Comment, &p.Book, &p.IdBook, &p.List, &p.Page, &p.DateAdded, &p.EnglishShort, &p.EnglishShortAuto)
	})
	return data
}

// GetWordsBySql Общая механика выборки из таблицы
func GetEnglishWordsBySql(database *sqlx.DB, sql string, scanner func(rows *sql.Rows, p *EnglishWord) error) []EnglishWord {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []EnglishWord
	for rows.Next() {
		p := EnglishWord{}
		err := scanner(rows, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
