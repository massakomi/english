package cmd

import (
	_ "database/sql"
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// с этим интерфейсом тут просто жесть. в любой момент может придти nil и все. везде нужны проверки. а код и так уже сумбурный с этим интерфейсом.

type assoc []map[string]interface{}

// getDataForHome данные для главной страницы, массив слов
func getDataForHome(context *gin.Context, database *sqlx.DB) assoc {
	limit := 5
	where := buildWhereForHomeData(context)
	//fmt.Println(where)
	sql := fmt.Sprintf(`select *
		from english_words
		where %v
		order by list desc, date_added desc
		limit %v`, where, limit)
	//fmt.Println(sql)
	data := db.GetData(sql, database)
	data, _ = fillShortAndCalcStat(data)
	data = prepareDataForTable(data, context, database)
	return data
}

// Подготовка данных перед выводом в таблицу
func prepareDataForTable(data assoc, context *gin.Context, database *sqlx.DB) assoc {
	baseStat := getBaseStat(data, database)
	//pretty.PrettyPrint(baseStat)

	cnt := len(data)
	for key, item := range data {

		if item["english_short_auto"] != nil {
			if baseStat[item["english_short_auto"].(string)] > 1 {
				item["englishBold"] = baseStat[item["english_short_auto"].(string)]
			}
		}

		format := "15:04"
		if context.Query("word") != "" {
			item["page"] = strconv.FormatInt(item["page"].(int64), 10) + " " + item["book"].(string)
			format = "Jan 02 15:04"
		}
		t := item["date_added"].(time.Time)
		item["date_added"] = t.Format(format)

		item["index"] = item["id"]
		if context.Query("book") != "" || context.PostForm("book") != "" {
			item["index"] = cnt - key
		}
	}
	return data
}

// Тут дозаполняется english_short, а вот для чего нужен stat без понятия
func fillShortAndCalcStat(data assoc) (assoc, map[string]int) {
	stat := map[string]int{}
	for _, value := range data {
		if value["english_short"] == nil {
			value["english_short"] = ""
		}
		stat[value["english"].(string)]++
		es := ""
		if value["english_short"].(string) == "" {
			es = text.BaseForm(value["english"].(string))
			if es != "" {
				value["english_short"] = es
			}
		} else {
			es = value["english_short"].(string)
		}
		if es != value["english"].(string) {
			stat[es]++
		}
	}
	return data, stat
}

// getBaseStat сколько слов в таблице имеют такие базовые формы
func getBaseStat(data assoc, database *sqlx.DB) map[string]int {
	baseStat := map[string]int{}
	if len(data) == 0 {
		return baseStat
	}
	bases := []string{}
	for _, item := range data {
		value := item["english_short_auto"].(string)
		if value != "" {
			bases = append(bases, value)
		}
	}
	sql := fmt.Sprintf(`select count(*) as c, english_short_auto
        from english_words
        where english_short_auto IN ('%v')
        group by 2`, strings.Join(bases, "','"))
	data = db.GetData(sql, database)
	for _, item := range data {
		baseStat[item["english_short_auto"].(string)] = int(item["c"].(int64))
	}
	return baseStat
}

// метод для getBaseStat
func buildWhereForHomeData(context *gin.Context) string {
	where := []string{"true"}
	if context.Query("word") != "" {
		word := context.Query("word")
		word = utils.PregReplace(word, "[\\s\\d]+$", "")
		word = text.BaseForm(word)
		where = append(where, fmt.Sprintf(`english_short_auto like '%v%%'`, word))
	}
	idBook := context.PostForm("book")
	if utils.IsNumeric(idBook) {
		where = append(where, fmt.Sprintf(`id_book = '%v'`, idBook))
	} else if context.Query("book") != "" {
		where = append(where, fmt.Sprintf(`book = '%v'`, context.Query("book")))
	}
	return strings.Join(where, " AND ")
}

// GetDataEnglishBooks Статистика чтения книг. Это массив структур EnglishBookread, где каждая структура - прочтение одной страницы книги
// Current true текущая страница читаю сейчас.
// Offset сколько секунд заняло прочтение этой страницы.
// DateAdded DateFinished - Даты начало конца чтение страницы
func GetDataEnglishBooks(max int, context *gin.Context, database *sqlx.DB) []models.EnglishBookread {
	dataEnglishAll := models.GetEnglishBookRead(database, context.Query("noexercises") == "", max)
	for key, item := range dataEnglishAll {
		if item.DateFinished == item.DateAdded {
			// echo notice := fmt.Sprintf(`Есть открытые страницы (%v, %v)`, item.name, item.page)
			item.Current = true
			item.DateFinished = time.Now()
		}
		item.Offset = item.DateFinished.Sub(item.DateAdded).Seconds()
		dataEnglishAll[key] = item
	}
	return dataEnglishAll
}

// ReadingStat блок статистики по чтению за определенный день
func ReadingStat(title string, dataEnglishBooks []models.EnglishBookread, date time.Time) string {
	ymd := date.Format("2006-01-02")
	inActive := true
	flow := map[string]string{}
	for key, item := range dataEnglishBooks {
		if key == 0 && !item.Current {
			inActive = false
		}
		added := item.DateAdded.Format("15:04")
		finished := item.DateFinished.Format("15:04")
		if item.DateAdded.Format("2006-01-02") != ymd {
			if len(flow) > 0 {
				break
			}
			continue
		}
		flow[added] = finished
	}

	// Получим ключи словаря в обратном порядке
	keys := slices.Sorted(maps.Keys(flow))
	slices.Reverse(keys)

	flowOut := map[string]string{}
	flowCounts := map[string]int{}
	root := ""
	finishPrev := ""
	for _, added := range keys {
		finish := flow[added]
		if IsEqualTimes(added, finishPrev) {
			flowOut[root] = finish
			flowCounts[root]++
		} else {
			root = added
			flowOut[added] = finish
			flowCounts[added]++
		}
		finishPrev = finish
	}

	var style string
	if inActive {
		style = ` style="color:green"`
	} else {
		style = ` style="color:#aaa"`
	}
	html := fmt.Sprintf(`<div%v>%v`, style, title)
	for key, value := range flowOut {
		if value == "" {
			value = "н.в."
		}
		html += fmt.Sprintf(` <span title='%v'>%v-%v &nbsp;</span>`, flowCounts[key], key, value)
	}
	html += `</div>`
	return html
}

// IsEqualTimes принимает строки вида "HH:mm" и возвращает true если между ними не более 3 минут
func IsEqualTimes(added string, finish string) bool {
	if finish == "" {
		return false
	}
	var hours, minutes time.Duration
	_, err := fmt.Sscanf(added, "%d:%d", &hours, &minutes)
	if err != nil {
		log.Fatalf(`Ошибка сканирования %v`, err)
	}
	t1 := minutes*time.Minute + hours*time.Hour

	_, err = fmt.Sscanf(finish, "%d:%d", &hours, &minutes)
	if err != nil {
		log.Fatalf(`Ошибка сканирования %v`, err)
	}
	t2 := minutes*time.Minute + hours*time.Hour

	return math.Abs(t1.Seconds()-t2.Seconds()) < 180
}

var BookStatus = map[string][]string{}
var BookStatusOnce sync.Once

// bookStyle стиль цвет книги в зависимости от статуса (прочитал или нет)
func bookStyle(bookName string, database *sqlx.DB) string {
	BookStatusOnce.Do(func() {
		books := models.GetBooksEx(database, "SELECT * FROM english_book")
		for _, item := range books {
			BookStatus[item.Status] = append(BookStatus[item.Status], item.Name)
		}
	})
	style := ""
	if BookStatus["read"] != nil && slices.Contains(BookStatus["read"], bookName) {
		style = "color:green"
	}
	if BookStatus["no"] != nil && slices.Contains(BookStatus["no"], bookName) {
		style = "color:red"
	}
	return style
}

// BooksByDays сгруппированные списки книг по дням
func BooksByDays(dataEnglishBooks []models.EnglishBookread, database *sqlx.DB) []map[string]any {
	bookStat, bookStatDmy, keys := statForBooksByDays(dataEnglishBooks)
	booksByDays := []map[string]any{}
	maxItems := 10
	for _, dmy := range keys {
		books := bookStat[dmy]
		list := []map[string]any{}
		totalPages := 0
		for bookName, item := range books {
			avg := item["seconds"] / item["pages"]
			list = append(list, map[string]any{
				"book":    bookName,
				"style":   bookStyle(bookName, database),
				"seconds": text.Fs(item["seconds"]),
				"pages":   item["pages"],
				"avg":     text.Fs(avg, true, false),
			})
			if bookName != "упражнения" {
				totalPages += item["pages"]
			}
		}
		booksByDays = append(booksByDays, map[string]any{
			"dmy":        dmy,
			"time":       text.Fs(bookStatDmy[dmy]),
			"totalPages": totalPages,
			"list":       list,
		})
		if maxItems--; maxItems == 0 {
			break
		}
	}
	return booksByDays
}

// statForBooksByDays собирает статистические словари для функции BooksByDays
func statForBooksByDays(dataEnglishBooks []models.EnglishBookread) (map[string]map[string]map[string]int, map[string]float64, []string) {
	bookStat := map[string]map[string]map[string]int{}
	bookStatDmy := map[string]float64{}
	for _, item := range dataEnglishBooks {
		date := item.DateAdded.Format("2006-01-02")
		if bookStat[date] == nil {
			bookStat[date] = map[string]map[string]int{}
		}
		if bookStat[date][item.Name] == nil {
			bookStat[date][item.Name] = map[string]int{}
		}
		bookStat[date][item.Name]["pages"]++
		bookStat[date][item.Name]["seconds"] += int(item.Offset)
		bookStatDmy[date] += item.Offset
	}
	keys := slices.Sorted(maps.Keys(bookStat))
	slices.Reverse(keys)
	return bookStat, bookStatDmy, keys
}

// Last5Pages для блока последних 5 страниц
func Last5Pages(dataEnglishBooks []models.EnglishBookread, context *gin.Context) []map[string]any {
	i := 0
	maxItems := 5
	data := []map[string]any{}
	for {
		item := dataEnglishBooks[i]
		i++
		if context.Query("book") != "" && context.Query("book") != item.Name {
			continue
		}
		data = append(data, map[string]any{
			"page":   item.Page,
			"name":   item.Name,
			"offset": text.Fs(item.Offset, true, false),
		})
		maxItems--
		if maxItems == 0 {
			break
		}
	}
	return data
}

// BookLast список книг, сортированные по дате последнего чтения с датой и количеством страниц
func BookLast(dataEnglishBooks []models.EnglishBookread, database *sqlx.DB) []map[string]any {
	bookLast := map[string]map[string]any{}
	bookForSort := map[string]int64{}
	for _, item := range dataEnglishBooks {
		if _, ok := bookLast[item.Name]; !ok {
			bookLast[item.Name] = map[string]any{
				"date":  item.DateAdded.Format("2006-01-02"),
				"page":  item.Page,
				"name":  item.Name,
				"style": bookStyle(item.Name, database),
			}
			bookForSort[item.Name] = item.DateAdded.Unix()
		}
	}
	keys := utils.MapKeySortByValues(bookForSort, true)
	var bookSlice []map[string]any
	for _, key := range keys {
		bookSlice = append(bookSlice, bookLast[key])
	}
	return bookSlice
}

// BookPages для роута /book список страниц одной книги
func BookPages(bookName string, dataEnglishBooks []models.EnglishBookread) []map[string]any {
	pages := make([]map[string]any, 0)
	last := ""
	for _, item := range dataEnglishBooks {
		if bookName == "" || bookName != item.Name {
			continue
		}
		date := item.DateAdded.Format("2006-01-02")
		info := map[string]any{
			"page":    item.Page,
			"offset":  text.Fs(item.Offset, true, false),
			"current": item.Current,
			"time":    item.DateAdded.Format("15:04"),
		}
		if date != last {
			info["date"] = date
		}
		pages = append(pages, info)
		last = item.DateAdded.Format("2006-01-02")
	}
	return pages
}
