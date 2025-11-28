package test

import (
	"english/cmd"
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobs/pretty"
	"log"
	"net/http"
	"time"
)

func TestGo() {
	//TestServer()
}

func TestFlow() {
	is := cmd.IsEqualTimes("12:40", "")
	fmt.Println(is)
}

func TestServer() {

	// gin.SetMode(gin.ReleaseMode)   debug off
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		database := db.Connect()
		defer database.Close()

		dataEnglishBooks := cmd.GetDataEnglishBooks(10, c, database)
		html := cmd.ReadingStat("Сегодня", dataEnglishBooks, time.Now())
		pretty.PrettyPrint(dataEnglishBooks)

		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"html": html,
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func TestBook() {
	database := db.Connect()
	//name := models.GetBookName(database, 1)
	//fmt.Println(name)
	models.AddBookPage(database, 1, 1, time.Now(), time.Now())
}

func TestFs() {

	base := text.Fs(200)
	fmt.Println(base)

	base = text.Fs(200, true, true)
	fmt.Println(base)
}

func TestBaseForm() {
	word := "upcoming"

	base := text.BaseForm(word)
	fmt.Println(base)

	base = text.BaseForm(word, true, true)
	fmt.Println(base)
}
