package cmd

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"strings"
)

func Run() {
	// gin.SetMode(gin.ReleaseMode)   debug off
	r := gin.Default()
	r.LoadHTMLGlob("public/*.html")
	r.HTMLRender = createMyRender()
	r.SetFuncMap(template.FuncMap{
		"lower":  strings.ToLower,
		"repeat": func(s string) string { return strings.Repeat(s, 2) },
	})
	r.Static("/static", "./public/static")

	r.GET("/", home)
	r.GET("/book", book)
	r.GET("/exercise", exercise)
	r.GET("/update-auto", updateAuto)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("home", "public/index.html", "public/home.html", "public/home.scripts.html", "public/home_table.html", "public/home_top.html")
	r.AddFromFiles("book", "public/index.html", "public/book.html")
	r.AddFromFiles("exercise", "public/index.html", "public/exercise/exercise.html", "public/exercise/exercise_scripts.html", "public/exercise/exercise_table.html")
	return r
}
