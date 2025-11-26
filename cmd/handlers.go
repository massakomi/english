package cmd

import (
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobs/pretty"
	"html/template"
	"net/http"
)

func home(c *gin.Context) {
	database := db.Connect()
	//fields := db.Fields(database, "english_book")
	//fmt.Println(fields)

	data := getDataForHome(c, database)
	pretty.PrettyPrint(data)
	baseStat := getBaseStat(data, database)
	pretty.PrettyPrint(baseStat)

	book := c.Query("book")

	books := models.GetBooks(database)

	c.HTML(http.StatusOK, "home.html", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, 0)),
		"book":             book,
	})
}

func updateAuto(c *gin.Context) {
	database := db.Connect()
	data := db.GetData(`select * from english_words where english_short_auto ='' or english_short_auto is null`, database)
	for _, item := range data {
		base := item["english_short"]
		if base == nil || base == "" {
			base = text.BaseForm(item["english"].(string))
			sql := fmt.Sprintf(`update english_words set english_short_auto='%v' where id=%v`, base, item["id"])
			database.MustExec(sql)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
