package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Book struct {
	Id         int
	Name       string
	Status     string
	Date_added string
}

func GetBooks(db *sql.DB) []Book {
	rows, err := db.Query("SELECT * FROM english_book ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	var books []Book
	defer rows.Close()
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.Id, &book.Name, &book.Status, &book.Date_added)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, book)
	}
	return books
}

func BooksSelector(books []Book, selected int) string {
	html := `<select required name="id_book" class="form-control form-control-sm me-2 w-50">`
	html += `<option value="">Книга</option>`
	for _, book := range books {
		add := ""
		if selected == book.Id {
			add = "selected"
		}
		html += fmt.Sprintf("<option %v value='%v'>%v</option>", add, book.Id, book.Name)
	}
	html += "</select>"
	return html
}
