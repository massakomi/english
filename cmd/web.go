package cmd

import (

	//"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"html/template"
	"log"
)

func Run() {
	// gin.SetMode(gin.ReleaseMode)   debug off
	r := gin.Default()
	addTemplates(r)
	r.Static("/static", "./public/static")

	r.GET("/", home)
	r.GET("/update-auto", updateAuto)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func addTemplates(router *gin.Engine) {
	files := []string{
		"./public/index.html",
		"./public/home.html",
		"./public/home.scripts.html",
	}
	html := template.Must(template.ParseFiles(files...))
	router.SetHTMLTemplate(html)
}
