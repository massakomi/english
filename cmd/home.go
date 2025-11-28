package cmd

import (
	_ "database/sql"
	"english/pkg/db"
	"english/pkg/models"
	"english/pkg/text"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"maps"
	"math"
	"slices"
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
func GetDataEnglishBooks(max int, c *gin.Context, database *sqlx.DB) []models.EnglishBookread {
	dataEnglishAll := models.GetEnglishBookRead(database, c.Query("noexercises") == "", max)
	for _, item := range dataEnglishAll {
		if item.DateFinished == item.DateAdded {
			// echo notice := fmt.Sprintf(`Есть открытые страницы (%v, %v)`, item.name, item.page)
			item.Current = true
			item.DateFinished = time.Now()
		}
		item.Offset = item.DateFinished.Sub(item.DateAdded)
	}
	return dataEnglishAll
}

func ReadingStat(title string, dataEnglishBooks []models.EnglishBookread, date time.Time) string {
	//ymd := date.Format("2006-01-02")
	inActive := true
	flow := map[string]string{}
	for key, item := range dataEnglishBooks {
		if key == 0 && !item.Current {
			inActive = false
		}
		added := item.DateAdded.Format("15:04")
		finished := item.DateFinished.Format("15:04")
		/*if item.DateAdded.Format("2006-01-02") != ymd {
			if len(flow) > 0 {
				break
			}
			continue
		}*/
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
