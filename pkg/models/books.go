package models

import (
	"english/pkg/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type EnglishBook struct {
	Id        int
	Name      string
	Status    string
	DateAdded time.Time
}

func GetBooksEx(database *sqlx.DB, sql string) []EnglishBook {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []EnglishBook
	for rows.Next() {
		p := EnglishBook{}
		err := rows.Scan(&p.Id, &p.Name, &p.Status, &p.DateAdded)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}

func GetBooks(database *sqlx.DB) []map[string]interface{} {
	return db.GetData("SELECT * FROM english_book ORDER BY name", database)
}

func BooksSelector(books []map[string]interface{}, selected int) string {
	html := `<select required name="id_book" class="form-control form-control-sm me-2 w-50">`
	html += `<option value="">Книга</option>`
	for _, book := range books {
		add := ""
		if selected == book["id"] {
			add = "selected"
		}
		html += fmt.Sprintf("<option %v value='%v'>%v</option>", add, book["id"], book["name"])
	}
	html += "</select>"
	return html
}

type BookRead struct {
	Name         string
	IdBook       int64 `db:"id_book"`
	Page         int
	DateAdded    time.Time
	DateFinished time.Time
}

func AddBookPage(database *sqlx.DB, id_book int64, page int, dateAdded time.Time, dateFinish time.Time) {
	name := GetBookName(database, id_book)
	_, err := database.NamedExec(""+
		`INSERT INTO english_bookread (name, id_book, page, date_added, date_finished) `+
		`VALUES (:name, :id_book, :page, :dateadded, :datefinished)`,
		&BookRead{name, id_book, page, dateAdded, dateFinish})
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(exec.LastInsertId())
	//fmt.Println(exec.RowsAffected())
}

func GetBookName(database *sqlx.DB, idBook int64) string {
	books := GetBooks(database)
	for _, item := range books {
		if item["id"].(int64) == idBook {
			return item["name"].(string)
		}
	}
	panic("book not found")
}

func DeleteBookPage(database *sqlx.DB, name string, page int) {
	database.MustExec("DELETE FROM english_bookread WHERE name=$1 AND page=$2", name, page)
}
