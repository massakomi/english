package models

import (
	"database/sql"
	"english/pkg/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"maps"
	"slices"
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

// GetEnglishBookReadBySql Общая механика выборки из таблиц книг в структуру (альтернативно выборке в словарь)
func GetEnglishBookReadBySql(database *sqlx.DB, sql string) []EnglishBookread {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	return scanAllEnglishBookRead(rows)
}

// GetEnglishBookReadByStmt Аналогичная функция, но принимает подготовленное выражение
func GetEnglishBookReadByStmt(database *sqlx.DB, sql string, data []any) []EnglishBookread {
	stmt, err := database.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(data...)
	return scanAllEnglishBookRead(rows)
}

func scanAllEnglishBookRead(rows *sql.Rows) []EnglishBookread {
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

// AveragePageTime среднее время прочтения страницы книги (выборка 10 страниц)
func AveragePageTime(database *sqlx.DB, bookId int) time.Duration {
	data := GetEnglishBookReadByStmt(
		database,
		`select * from english_bookread where id_book = $1 AND date_finished is not null order by id desc limit 10`,
		[]any{bookId},
	)
	var averagePageTime time.Duration
	for _, item := range data {
		averagePageTime += item.DateFinished.Sub(item.DateAdded)
	}
	averagePageTime = averagePageTime / time.Duration(len(data))
	return averagePageTime
}

// AutoPagination автоматически добавляет статистику прочитанных страниц по указанной книге с последней страницы в базе до страницы readpage
// ставит время прочтения среднее на основе функции AveragePageTime
func AutoPagination(book int, readpage int, database *sqlx.DB) {
	ymd := time.Now().Format("2006-01-02")
	sql := fmt.Sprintf(`select MIN(date_added) as dt from english_bookread where TO_CHAR(date_added, 'YYYY-MM-DD') = '%v'`, ymd)
	//sql := `select MIN(date_added) as dt from english_bookread`
	t := db.GetFirstVal(sql, database)
	var timeFinish time.Time
	if t == "" {
		timeFinish = time.Now()
	} else {
		var err error
		timeFinish, err = time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			panic(err)
		}
	}

	read := GetLastBookPage(database, book)
	tms := AveragePageTime(database, book)
	if readpage <= read.Page {
		panic("readpage <= read.Page")
	}
	timeAdded := timeFinish.Add(-tms)
	data := make(map[int]map[string]time.Time)
	for page := readpage; page > read.Page; page-- {
		if data[page] == nil {
			data[page] = map[string]time.Time{}
		}
		data[page]["timeAdded"] = timeAdded
		data[page]["timeFinish"] = timeFinish
		timeFinish = timeAdded
		timeAdded = timeAdded.Add(-tms)
	}

	pages := slices.Sorted(maps.Keys(data))
	slices.Reverse(pages)

	for _, page := range pages {
		AddBookPage(database, int64(book), page, data[page]["timeAdded"], data[page]["timeFinish"])
	}

}

// GetLastBookPage
func GetLastBookPage(database *sqlx.DB, bookId int) EnglishBookread {
	data := GetEnglishBookReadByStmt(
		database,
		`select * from english_bookread where id_book=$1 order by id desc limit 1`,
		[]any{bookId},
	)
	return data[0]
}
