package cmd

import (
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/gobs/pretty"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func home(context *gin.Context) {
	database := db.Connect()
	data := getDataForHome(context, database)
	books := models.GetBooks(database)

	maxDays, _ := strconv.Atoi(context.DefaultQuery("max", "10"))
	dataEnglishBooks := GetDataEnglishBooks(maxDays, context, database)
	//pretty.PrettyPrint(dataEnglishBooks)

	html := ReadingStat("Сегодня", dataEnglishBooks, time.Now())
	html += ReadingStat("Вчера", dataEnglishBooks, time.Now().Add(-time.Hour*24))

	context.HTML(http.StatusOK, "home", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, utils.GetPostDefaultInt("book", context))),
		"book":             context.Query("book"),
		"word":             context.Query("word"),
		"cookieBook":       utils.GetCookie("book", context),
		"data":             data,
		"readingStat":      template.HTML(html),
		"bookLast":         BookLast(dataEnglishBooks, database),
		"booksByDays":      BooksByDays(dataEnglishBooks, database),
		"last5pages":       Last5Pages(dataEnglishBooks, context),
		"template":         "home",
	})
}

func book(context *gin.Context) {
	database := db.Connect()
	books := models.GetBooks(database)
	maxDays, _ := strconv.Atoi(context.DefaultQuery("max", "10"))
	dataEnglishBooks := GetDataEnglishBooks(maxDays, context, database)

	context.HTML(http.StatusOK, "book", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, utils.GetPostDefaultInt("book", context))),
		"book":             context.Query("book"),
		"bookPages":        BookPages(context.Query("book"), dataEnglishBooks),
		"template":         "book",
	})
}

func exercise(context *gin.Context) {
	database := db.Connect()
	books := models.GetBooks(database)

	context.HTML(http.StatusOK, "exercise", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, utils.GetPostDefaultInt("book", context))),
		"data":             GetDataForList(database),
	})
}

func exercisePage(context *gin.Context) {
	database := db.Connect()
	books := models.GetBooks(database)
	index := context.Param("index")
	outputData, exerciseComment := GetDataForExercise(database, index)

	context.HTML(http.StatusOK, "exercise", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, utils.GetPostDefaultInt("book", context))),
		"exerciseIndex":    index,
		"exerciseDbUrl":    `http://msc/index.php?db=tester&table=english_exercise_questions&s=tbl_data&where=exercise=` + index,
		"date":             time.Now().Format("15:04:05"),
		"exerciseComment":  exerciseComment,
		"data":             outputData,
		"exerciseInfo":     GetExerciseStarted(database, index),
	})
}

func exerciseStart(context *gin.Context) {
	database := db.Connect()
	index := context.Param("index")
	UpdateExerciseIfStarted(database, index, true)
	context.String(http.StatusOK, "")
}

func exerciseRegister(context *gin.Context) {
	exercise := context.Param("index")
	tm := context.Query("time")
	question := context.Query("index")
	errors := context.Query("errors")
	comment := context.Query("comment")
	if exercise == "" || question == "" || question == "undefined" {
		context.String(http.StatusBadRequest, "")
	}
	database := db.Connect()
	now := time.Now().Add(-time.Hour).Format("2006-01-02 15:04:05")
	where := fmt.Sprintf(`exercise=%v AND question='%v' AND date_added > '%v'`, exercise, question, now)
	item := models.GetExerciseQuestionsByWhere(database, where)

	data := make(map[string]any)
	data["date_added"] = time.Now().Format("2006-01-02 15:04:05")
	if utils.IsNumeric(errors) {
		data["errors"] = errors
	}
	if comment != "" {
		data["comment"] = comment
	}
	if tm != "" {
		data["time"] = tm
	}

	if len(item) > 0 {
		db.Update(database, "english_exercise_questions", item[0].Id, data)
	} else {
		data["exercise"] = exercise
		data["question"] = question
		db.Insert(database, "english_exercise_questions", data)
	}

	UpdateExerciseIfStarted(database, exercise, true)

	context.String(http.StatusOK, "ok")
}

func updateAuto(context *gin.Context) {
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
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
