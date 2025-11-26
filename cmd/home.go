package cmd

import (
	"english/pkg/db"
	"english/pkg/text"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"strings"
)

// getBaseStat сколько слов в таблице имеют такие базовые формы
func getBaseStat(data []map[string]interface{}, database *sqlx.DB) map[string]int {
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

// getDataForHome данные для главной страницы, массив слов
func getDataForHome(c *gin.Context, database *sqlx.DB) []map[string]interface{} {
	limit := 2
	where := buildWhereForHomeData(c)
	sql := fmt.Sprintf(`select *
		from english_words
		where %v
		order by list desc, date_added desc
		limit %v`, where, limit)
	data := db.GetData(sql, database)
	return data
}

func buildWhereForHomeData(c *gin.Context) string {
	where := []string{"true"}
	if c.Query("word") != "" {
		word := c.Query("word")
		word = utils.PregReplace(word, "[\\s\\d]+$", "")
		word = text.BaseForm(word)
		where = append(where, fmt.Sprintf(`english_short_auto like '%v%'`, word))
	}
	idBook := c.PostForm("book")
	if utils.IsNumeric(idBook) {
		where = append(where, fmt.Sprintf(`id_book = '%v'`, idBook))
	} else if c.Query("book") != "" {
		where = append(where, fmt.Sprintf(`book = '%v'`, c.Query("book")))
	}
	return strings.Join(where, " AND ")
}
