package cmd

import (
	"english/pkg/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func GetDataForArticles(articlesMode bool) {
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

	defer timer("main")()
	chs := utils.PregSplit(`(?s)[\.?!]\s*`, text)
	for _, item := range chs {
		fmt.Println(item)
	}
	fmt.Println(len(chs))
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
