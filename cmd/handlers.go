package cmd

import (
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"

	//"github.com/gobs/pretty"
	"html/template"
	"net/http"
	"time"
)

func home(c *gin.Context) {
	database := db.Connect()
	data := getDataForHome(c, database)
	//pretty.PrettyPrint(data)
	books := models.GetBooks(database)

	maxDays, _ := strconv.Atoi(c.DefaultQuery("max", "10"))
	dataEnglishBooks := GetDataEnglishBooks(maxDays, c, database)
	// pretty.PrettyPrint(dataEnglishBooks)

	html := ReadingStat("Сегодня", dataEnglishBooks, time.Now())
	html += ReadingStat("Вчера", dataEnglishBooks, time.Now().Add(-time.Hour*24))

	c.HTML(http.StatusOK, "home.html", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, 0)),
		"book":             c.Query("book"),
		"word":             c.Query("word"),
		"cookieBook":       utils.GetCookie("book", c),
		"countData":        len(data),
		"data":             data,
		"readingStat":      template.HTML(html),
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
