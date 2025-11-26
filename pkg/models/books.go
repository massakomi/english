package models

import (
	"english/pkg/db"
	"fmt"
	"github.com/jmoiron/sqlx"
)

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
