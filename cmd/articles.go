package cmd

import (
	"english/pkg/models"
	"english/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetDataForArticles(context *gin.Context, articlesMode bool) (string, int, int) {
	content, err := os.ReadFile("data/exercise-text.txt") // data/Fox-Street.txt
	if err != nil {
		log.Fatal(err)
	}
	text := strings.Replace(string(content), `.”`, `”.`, -1)

	if articlesMode {
		articlesPrepare(&text)
	} else {
		predlogsPrepare(&text)
	}

	selected := 1
	nextPage, _ := context.Cookie("article-next-page")
	if nextPage != "" {
		selected, _ = strconv.Atoi(nextPage)
	}

	chs := utils.PregSplit(`(?s)[\.?!]\s*`, text)

	articles := ""
	chunk := ""
	counts := 0
	for _, item := range chs {
		chunk += strings.TrimSpace(item) + "."
		l := len(utils.StripTags(chunk))
		if l > 200 {
			addChunk(&articles, selected, &chunk, &counts)
		}
	}
	if chunk != "" {
		addChunk(&articles, selected, &chunk, &counts)
	}

	return articles, selected, counts
}

func addChunk(articles *string, selected int, chunk *string, counts *int) {
	*counts++
	add := ` class="hidden"`
	if *counts == selected {
		add = ""
	}
	*articles += fmt.Sprintf(`<p%v>%v</p>`, add, *chunk)
	*chunk = ""
}

func predlogsPrepare(text *string) {
	//extra := `|above|below|over|under|before|behind|among|between|by|near|beside|next to|beyond|across|opposite|in front of|inside|outside|from|towards|across|through`
	predlogsList := []string{
		`on|in|at|of|to|into|for|from|by`,
		`with|near|onto|upon`,
	}
	*text = utils.PregReplace(*text, `(?i)((had|have|tried|begun|nothing|sorry|supposed|going|want|able) to|(instead|in order|one) of)`, `$1*`)
	for _, predlogs := range predlogsList {
		a := strings.Split(predlogs, "|")
		for key, item := range a {
			a[key] = fmt.Sprintf(`<i>%v</i>`, item)
		}
		predlogsHtml := strings.Join(a, `<s>|</s>`)
		pattern := `(?i)[\s]('` + predlogs + `')([\s,\.?!:;”])`
		replace := fmt.Sprintf(` <span data-title="$1" class="predlog">%v</span>$2`, predlogsHtml)
		*text = utils.PregReplace(*text, pattern, replace)
	}
	*text = strings.Replace(*text, `*`, ``, -1)
}

func articlesPrepare(text *string) {
	*text = utils.PregReplace(*text, `(?s)[\r\n]+`, " ")
	*text = utils.PregReplaceCallback(*text, `(The|A|An) ([a-z])`, func(str string) string {
		a := strings.Split(str, " ")
		return a[0] + " " + strings.ToTitle(a[1])
	})
	*text = strings.Replace(*text, `Mr.`, `Mr`, -1)
	*text = strings.Replace(*text, `Mrs.`, `Mrs`, -1)
	*text = utils.PregReplace(*text, `(?i)([\s“"]|^)(the) ([^\s]+)`, `$1<b>$2</b> <span data-title="$2">$3</span>`)
	*text = utils.PregReplace(*text, `(?i)([\s“"]|^)(an) ([^\s]+)`, `$1<b>$2</b> <span data-title="$2">$3</span>`)
}

// Для страницы memory

func GetDataForMemory(database *sqlx.DB) []map[string]any {
	output := make([]map[string]any, 0)
	data, statData := GetWords(database)
	for _, item := range data {
		stat := statData[item.Id]
		add := ""
		if stat != nil {
			result, _ := strconv.Atoi(fmt.Sprint(stat["result"]))
			if result > 0 {
				add = ` class="text-ok"`
			} else {
				add = ` ' class="text-danger"'`
			}
		}

		output = append(output, map[string]any{
			"id":           item.Id,
			"result":       stat["result"],
			"add":          template.HTMLAttr(add),
			"word":         item.Word,
			"translate":    item.Translate,
			"remindOffset": getRemindOffset(item),
		})
	}
	return output
}

func getRemindOffset(item models.Word) int {
	var remindOffset time.Duration
	if item.DateRemind.Valid {
		remindOffset = time.Now().Sub(item.DateRemind.Time)
	}
	remindSeconds := int(remindOffset.Seconds())
	if remindOffset.Seconds() > 3600*12 {
		remindSeconds = 0
	}
	return remindSeconds
}

func GetWords(database *sqlx.DB) ([]models.Word, map[int]map[string]any) {
	data := models.GetWordsWithDateResult(database)
	a := models.GetLogs(database)
	countWords := make(map[int]int)
	statData := make(map[int]map[string]any)
	for _, item := range a {
		if statData[item.IdWord] == nil {
			statData[item.IdWord] = make(map[string]any)
		}
		if statData[item.IdWord]["result"] == nil {
			statData[item.IdWord]["result"] = item.Result
			statData[item.IdWord]["date"] = item.DateAdded
			statData[item.IdWord]["time"] = time.Now().Sub(item.DateAdded)
		}
		countWords[item.IdWord]++
		statData[item.IdWord]["count"] = countWords[item.IdWord]
	}
	return data, statData
}
