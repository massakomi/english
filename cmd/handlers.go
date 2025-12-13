package cmd

import (
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"

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
		"hideAfterStart":   "true",
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
	time := context.Query("time")
	question := context.Query("index")
	errors := context.Query("errors")
	comment := context.Query("comment")
	if exercise == "" || question == "" || question == "undefined" {
		context.String(http.StatusBadRequest, "")
	}
	database := db.Connect()
	ExerciseAddOrUpdate(database, exercise, time, question, errors, comment)
	context.String(http.StatusOK, "ok")
}

func exerciseArticles(context *gin.Context) {
	database := db.Connect()
	books := models.GetBooks(database)

	articlesMode := strings.Contains(context.FullPath(), "articles")
	exercise := "articles"
	if !articlesMode {
		exercise = "predlogs"
	}
	articles, selected, counts := GetDataForArticles(context, articlesMode)

	context.HTML(http.StatusOK, "exercise", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, utils.GetPostDefaultInt("book", context))),
		"exerciseIndex":    exercise,
		"exerciseInfo":     GetExerciseStarted(database, exercise),
		"date":             time.Now().Format("15:04:05"),
		"hideAfterStart":   "false",
		"articles":         template.HTML(articles),
		"selected":         selected,
		"counts":           counts,
		"articlesMode":     articlesMode,
	})
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

func memory(context *gin.Context) {
	database := db.Connect()
	books := models.GetBooks(database)

	context.HTML(http.StatusOK, "memory", gin.H{
		"getBooksSelector": template.HTML(models.BooksSelector(books, utils.GetPostDefaultInt("book", context))),
		"data":             GetDataForMemory(database),
	})
}

func wordLog(context *gin.Context) {
	database := db.Connect()
	db.Insert(database, "english_log", map[string]any{
		"id_word":    context.Param("id"),
		"date_added": time.Now().Format("2006-01-02 15:04:05"),
		"result":     context.Param("rate"),
	})
	context.String(http.StatusOK, "")
}

// Загружает на редактирование слово в /memory
func wordEdit(context *gin.Context) {
	database := db.Connect()
	words := models.GetWords(database, fmt.Sprintf(`id=%v`, context.Param("id")))
	context.JSON(http.StatusOK, gin.H{
		"id":        words[0].Id,
		"word":      words[0].Word,
		"translate": words[0].Translate,
		"comment":   words[0].Comment.String,
	})
}

// Сохраняет слово в /memory
func wordEditSave(context *gin.Context) {
	database := db.Connect()
	id, _ := strconv.Atoi(context.Param("id"))
	db.Update(database, "english", id, map[string]any{
		"word":      context.PostForm("word"),
		"translate": context.PostForm("translate"),
		"comment":   context.PostForm("comment"),
	})
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func wordSave(context *gin.Context) {
	data := make(map[string]any)
	id, _ := strconv.Atoi(context.Param("id"))
	if russian, _ := context.GetQuery("russian"); russian != "" {
		data["russian"] = russian
	}
	if english, _ := context.GetQuery("english"); english != "" {
		data["english"] = english
	}
	database := db.Connect()
	db.Update(database, `english_words`, id, data)
	context.JSON(http.StatusOK, gin.H{
		"id":      id,
		"russian": data["russian"],
		"english": data["english"],
	})
}

// todo $.get('?action=get_word&word='+this.value,
func wordGet(context *gin.Context) {
	database := db.Connect()
	cookBook, _ := context.Cookie("book")
	word := context.Query("word")
	data := models.GetEnglishWords(database, fmt.Sprintf(`"english" like '%v' order by id desc limit 5`, word))
	if len(data) == 0 {
		data = models.GetEnglishWords(database, fmt.Sprintf(`"english" like '%v%%' order by id desc limit 5`, word))
	} else if cookBook != "" {
		var ids []string
		for _, item := range data {
			ids = append(ids, strconv.FormatInt(int64(item.Id), 10))
		}
		sql := fmt.Sprintf(`"english" like '%v%%' AND id_book='%v' AND id NOT IN (%v) order by id desc`, word, cookBook, strings.Join(ids, ","))
		data2 := models.GetEnglishWords(database, sql)
		if len(data2) > 0 {
			data = append(data, data2...)
		}
	}
	if len(data) == 0 {
		word := text.BaseForm(context.PostForm("word"), false)
		where := fmt.Sprintf(`"english" like '%v%%'`, word)
		if strings.Contains(word, " ") {
			where += fmt.Sprintf(` OR "english" like '%v%%'`, strings.Replace(word, " ", "-", -1))
		}
		if strings.Contains(word, "-") {
			where += fmt.Sprintf(` OR "english" like '%v%%'`, strings.Replace(word, "-", " ", -1))
		}
		data = models.GetEnglishWords(database, where)
	}
	output := []string{}
	for _, item := range data {
		book := item.Book
		if cookBook != "" && cookBook == strconv.FormatInt(int64(item.IdBook), 10) {
			book = fmt.Sprintf(`<b>%v</b>`, book)
		}
		itemPage, _ := item.Page.Value()
		output = append(output, fmt.Sprintf(`%v (%v - %v)`, item.Russian, book, itemPage))
	}
	context.String(http.StatusOK, strings.Join(output, " / "))
}

func wordTable(context *gin.Context) {
	database := db.Connect()
	data := getDataForHome(context, database)
	context.HTML(http.StatusOK, "home_table", gin.H{
		"data": data,
	})
}

func bookPageGet(context *gin.Context) {
	bookId := context.Param("id")
	database := db.Connect()
	sql := fmt.Sprintf(`select MAX(page) as dt from english_bookread where id_book=%v`, bookId)
	c := db.GetFirstVal(sql, database)
	context.JSON(http.StatusOK, gin.H{
		"page": c,
	})
}

func translateAdd(context *gin.Context) {
	database := db.Connect()
	base := text.BaseForm(context.PostForm("english"))
	id := 1
	idBook := utils.GetCookie("book", context)
	if idBook != "" {
		id, _ = strconv.Atoi(idBook)
	}
	data := map[string]any{
		"english":            strings.TrimSpace(context.PostForm("english")),
		"russian":            strings.TrimSpace(context.PostForm("russian")),
		"book":               models.GetBookName(database, int64(id)),
		"id_book":            id,
		"list":               20,
		"date_added":         time.Now().Format("2006-01-02 15:04:05"),
		"english_short_auto": base,
		"page":               utils.GetCookie("book", context),
	}
	res := db.Insert(database, "english_words", data)
	idx, _ := res.LastInsertId()
	rows, _ := res.RowsAffected()
	context.JSON(http.StatusOK, gin.H{
		"id":   idx,
		"rows": rows,
	})
}

func bookRead(context *gin.Context) {
	database := db.Connect()
	idBook := context.PostForm("id_book")
	idBookInt, _ := strconv.Atoi(idBook)
	utils.SetCookie("book", idBook, context, 24*365)

	a := models.GetLastBookPage(database, idBookInt)
	currentPage := 0
	if a.DateFinished == a.DateAdded {
		currentPage = a.Page
	}

	readPage := context.PostForm("readpage")
	readPageInt, _ := strconv.Atoi(readPage)

	fillPages := context.PostForm("fill-pages")
	if fillPages != "" {
		models.AutoPagination(idBookInt, readPageInt, database)
		context.JSON(http.StatusOK, gin.H{
			"ok": "fill-pages!",
		})
	}

	finish := readPageInt
	if currentPage != 0 && currentPage != readPageInt {
		finish = currentPage
	}

	// Добавление промежуточных страниц для случаев "перескока" страницы через 2 и больше
	addNewPage := addIntermediatePages()

	if finish > 0 {
		// todo
	}

	utils.SetCookie("page", readPage, context, 1)
	if addNewPage {
		models.AddBookPage(database, int64(idBookInt), readPageInt, time.Now(), time.Now())
	}

	context.JSON(http.StatusOK, gin.H{
		"ok": 1,
	})
}

func addIntermediatePages() bool {
	addNewPage := true
	// todo
	return addNewPage
}
