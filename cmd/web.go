package cmd

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func Run() {
	r := gin.Default()
	r.LoadHTMLGlob("public/*.html")
	r.HTMLRender = createMyRender()
	r.Static("/static", "./public/static")

	r.GET("/", home)
	r.POST("/translate/add", translateAdd)
	r.POST("/book/read", bookRead)
	r.GET("/book", book)
	r.GET("/exercise", exercise)
	r.GET("/exercise/:index", exercisePage)
	r.GET("/exercise/start/:index", exerciseStart)
	r.GET("/exercise/register/:index", exerciseRegister)
	r.GET("/exercise/articles", exerciseArticles)
	r.GET("/exercise/prepositions", exerciseArticles)
	r.GET("/update-auto", updateAuto)
	r.GET("/memory", memory)
	r.GET("/word/log/:id/:rate", wordLog)
	r.GET("/word/edit/:id", wordEdit)
	r.POST("/word/edit/:id", wordEditSave)
	r.GET("/word/save/:id", wordSave)
	r.GET("/word/get/:id", wordGet)
	r.GET("/word/data/:id", wordData)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("home", "public/index.html", "public/home.html", "public/home.scripts.html", "public/home_table.html", "public/home_top.html")
	r.AddFromFiles("book", "public/index.html", "public/book.html")
	r.AddFromFiles("memory", "public/index.html", "public/memory.html")
	r.AddFromFiles(
		"exercise",
		"public/index.html",
		"public/exercise/exercise.html",
		"public/exercise/exercise_scripts.html",
		"public/exercise/exercise_table.html",
		"public/exercise/questions.html",
		"public/exercise/articles.html",
	)
	return r
}
