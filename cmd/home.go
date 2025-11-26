package cmd

import (
	_ "database/sql"
	"english/pkg/db"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"

	//"github.com/gobs/pretty"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

// с этим интерфейсом тут просто жесть. в любой момент может придти nil и все. везде нужны проверки. а код и так уже сумбурный с этим интерфейсом.

type assoc []map[string]interface{}

// getDataForHome данные для главной страницы, массив слов
func getDataForHome(c *gin.Context, database *sqlx.DB) assoc {
	limit := 2
	where := buildWhereForHomeData(c)
	//fmt.Println(where)
	sql := fmt.Sprintf(`select *
		from english_words
		where %v
		order by list desc, date_added desc
		limit %v`, where, limit)
	//fmt.Println(sql)
	data := db.GetData(sql, database)
	data, _ = fillShortAndCalcStat(data)
	data = prepareDataForTable(data, c, database)
	return data
}

// Подготовка данных перед выводом в таблицу
func prepareDataForTable(data assoc, c *gin.Context, database *sqlx.DB) assoc {
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
		if c.Query("word") != "" {
			item["page"] = strconv.FormatInt(item["page"].(int64), 10) + " " + item["book"].(string)
			format = "Jan 02 15:04"
		}
		t := item["date_added"].(time.Time)
		item["date_added"] = t.Format(format)

		item["index"] = item["id"]
		if c.Query("book") != "" || c.PostForm("book") != "" {
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
func buildWhereForHomeData(c *gin.Context) string {
	where := []string{"true"}
	if c.Query("word") != "" {
		word := c.Query("word")
		word = utils.PregReplace(word, "[\\s\\d]+$", "")
		word = text.BaseForm(word)
		where = append(where, fmt.Sprintf(`english_short_auto like '%v%%'`, word))
	}
	idBook := c.PostForm("book")
	if utils.IsNumeric(idBook) {
		where = append(where, fmt.Sprintf(`id_book = '%v'`, idBook))
	} else if c.Query("book") != "" {
		where = append(where, fmt.Sprintf(`book = '%v'`, c.Query("book")))
	}
	return strings.Join(where, " AND ")
}

// Статистика чтения книг
func getDataEnglishBooks(max int, c *gin.Context, database *sqlx.DB) []EnglishExercise {
	add := ""
	if c.Query("noexercises") == "" {
		add = fmt.Sprintf(`UNION (
			select id, 'упражнения' as name, 0 AS id_book, page, date_added, date_finished 
			from english_exercise 
			where date_added > date_add(NOW(), -'%v day'::interval))`, max)
	}
	sql := fmt.Sprintf(`select * from english_bookread %v order by date_added desc limit 3`, add)
	//fmt.Println(sql)
	dataEnglishAll := getDataEnglish(sql, database)

	for _, item := range dataEnglishAll {
		if item.dateFinished == item.dateAdded {
			// echo notice := fmt.Sprintf(`Есть открытые страницы (%v, %v)`, item.name, item.page)
			item.current = true
			item.dateFinished = time.Now()
		}
		item.offset = item.dateFinished.Sub(item.dateAdded)
	}

	return dataEnglishAll
}

type EnglishExercise struct {
	id           int
	name         string
	idBook       int
	page         int
	dateAdded    time.Time
	dateFinished time.Time
	current      bool
	offset       time.Duration
}

// Общая механика выборки из таблиц книг в структуру (альтернативно выборке в словарь)
func getDataEnglish(sql string, db *sqlx.DB) []EnglishExercise {
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items []EnglishExercise
	for rows.Next() {
		p := EnglishExercise{}
		err := rows.Scan(&p.id, &p.name, &p.idBook, &p.page, &p.dateAdded, &p.dateFinished)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, p)
	}
	return items
}
