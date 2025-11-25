package cmd

import (
	"english/pkg/db"
	"english/pkg/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func home(c *gin.Context) {
	database := db.Connect()
	//fields := db.Fields(database, "english_book")
	//fmt.Println(fields)

	books := models.GetBooks(database)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, 0)),
	})
}
