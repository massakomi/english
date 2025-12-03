package cmd

import (
	"english/pkg/models"
	"english/pkg/utils"
	"fmt"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func addExercise() {
}

func updateExercise() {
}

func updateExerciseIfStarted() {
}

func getExerciseStarted() {
}

func ExerciseStat(database *sqlx.DB) map[int]models.Exercise {
	stat := map[int]models.Exercise{}
	data := models.GetExercises(database)
	for _, item := range data {
		if stat[item.Page].Page > 0 {
			continue
		}
		dayAgo := time.Now().Sub(item.DateAdded)
		item.DaysAgo = int(math.Round(dayAgo.Seconds() / 86400))
		t := item.DateFinished.Sub(item.DateAdded)
		item.Time = utils.RoundFloat(t.Seconds()/60, 1)
		stat[item.Page] = item
	}
	return stat
}

func getExercisesQuestionsStat(database *sqlx.DB) map[int][]int {
	data := models.GetExerciseQuestions(database)
	cache := make(map[int]map[string]int)
	stats := make(map[int][]int)
	for _, item := range data {
		if cache[item.Exercise][item.Question] > 0 {
			continue
		}
		if cache[item.Exercise] == nil {
			cache[item.Exercise] = make(map[string]int)
		}
		cache[item.Exercise][item.Question] = item.Errors
		stats[item.Exercise] = append(stats[item.Exercise], item.Errors)
	}
	return stats
}

func GetTotalErrorsByExercises(database *sqlx.DB) map[int]map[string]int {
	errors := make(map[int]map[string]int)
	stats := getExercisesQuestionsStat(database)
	for key, item := range stats {
		errors[key] = map[string]int{
			"errors": utils.SumIntSlice(item),
			"count":  len(item),
		}
	}
	return errors
}

var totalStat = make(map[int]models.Exercise)
var totalErrors = make(map[int]map[string]int)

func GetExerciseStatStyle(database *sqlx.DB, index int) models.Exercise {
	statItem := models.Exercise{}
	if len(totalStat) == 0 {
		totalStat = ExerciseStat(database)
		totalErrors = GetTotalErrorsByExercises(database)
	}
	style := `text-decoration:none; `
	if totalStat[index].Page > 0 {
		statItem = totalStat[index]
		if statItem.DaysAgo < 7 {
			style += `font-weight:bold; color:green; `
			if statItem.DaysAgo < 1 {
				style += `font-size:13px; `
			}
		} else if statItem.DaysAgo < 14 {
			style += `color:green; `
		} else if statItem.DaysAgo < 30 {
			style += `color:#aaa; `
		} else {
			style += `color:#ccc; `
		}
		statItem.Style = template.CSS(style)
		statItem.Errors = totalErrors[index]["errors"]
	}
	return statItem
}

func formatOffset(offsetSeconds float64) string {
	offsetInfo := ""
	offset := utils.RoundFloat(offsetSeconds/86400, 1)
	if offsetSeconds > 86400 {
		offsetInfo = fmt.Sprintf("%v days", offset)
	} else if offsetSeconds > 3600 {
		offset = utils.RoundFloat(offsetSeconds/3600, 1)
		offsetInfo = fmt.Sprintf("%v hours", offset)
	} else {
		offset = utils.RoundFloat(offsetSeconds/60, 1)
		offsetInfo = fmt.Sprintf("%v mins", offset)
	}
	return offsetInfo
}

func GetExerciseQuestion(database *sqlx.DB, exercise string, last bool) map[int]map[string]map[string]any {
	questions := make(map[int]map[string]map[string]any)
	data := models.GetExerciseQuestionsByWhere(database, "exercise="+exercise)
	for _, item := range data {
		tm, _ := time.Parse("2006-01-02 15:04:05", item.DateAdded)
		dmy := tm.Format("2006-01-02")
		offsetSeconds := time.Now().Sub(tm).Seconds()
		offsetDays := int(math.Floor(offsetSeconds / 86400))
		offsetInfo := formatOffset(offsetSeconds)
		row := map[string]any{
			"errors":      item.Errors,
			"time":        "",
			"offset":      offsetDays,
			"offset-info": offsetInfo,
		}
		if item.Time.Valid {
			row["time"] = fmt.Sprintf(`%vs`, item.Time.Int64)
		}
		question, _ := strconv.Atoi(item.Question)
		if last {
			if questions[question] != nil {
				continue
			}
		}
		if questions[question] == nil {
			questions[question] = make(map[string]map[string]any)
		}
		questions[question][dmy] = row
	}
	return questions
}

func GetExerciseQuestionLast(database *sqlx.DB, exercise string) map[int]map[string]any {
	outputData := GetExerciseQuestion(database, exercise, true)
	questions := make(map[int]map[string]any)
	for question, values := range outputData {
		for _, valuesInner := range values {
			questions[question] = valuesInner
		}
	}
	return questions
}

var exerciseStat map[int]map[string]map[string]any

func exerciseQuestionStat(database *sqlx.DB, exercise string, index int) string {
	if exerciseStat == nil {
		exerciseStat = GetExerciseQuestion(database, exercise, false)
	}
	add := ""
	if exerciseStat[index] != nil {
		qs := []string{}
		for _, item := range exerciseStat[index] {
			e := fmt.Sprintf(`<span style="color:red">%v errors</span>`, item["errors"])
			if item["time"] != " " {
				e = fmt.Sprintf(`%v `, item["time"]) + e
			}
			qs = append(qs, e)
		}
		slices.Reverse(qs)
		add = fmt.Sprintf(` <span style="font-size:11px; color:green">%v</span>`, strings.Join(qs, " / "))
	}
	return add
}

func getAllExercises() [][]string {
	content, err := os.ReadFile(`data/exercise.txt`)
	if err != nil {
		log.Fatal(err)
	}
	data := utils.PregMatchAllEx(`(?is)\s*Еxеrcіsе (\d+)[\s\.]+[^\r\n]+`, string(content))
	return data
}

func extractSentences(rus string) []string {
	rus = strings.TrimSpace(rus)
	a := utils.PregSplit(`(?i)(^|\s)(\d+)[,\.](\s|$)`, rus)
	a = slices.DeleteFunc(a, func(n string) bool {
		return n == ""
	})
	return a
}

func ExtractContentSentences(content string) ([]string, []string) {
	contentOrig := content
	content = strings.TrimSpace(content)
	if utils.PregMatch(`(?i)Еxеrcіsе (\d+)`, content) != "" {
		content = utils.PregReplace(content, `(?i)Еxеrcіsе (\d+)[\s\.]+`, "")
	}
	content = strings.Replace(content, "\r\n", "\n", -1)
	content = utils.PregReplace(content, `(?i)([\d])t`, "$1.")
	content = utils.PregReplace(content, `(?s)\n[ \t]+`, "\n")
	data := utils.PregSplit("(?s)[\r\n]{2,}", content)
	if len(data) < 2 {
		fmt.Println(contentOrig)
		log.Fatal("Не удалось сплитить", content)
	}
	rus := extractSentences(data[0])
	eng := extractSentences(data[1])
	return rus, eng
}

func GetAllTexts() map[string]string {
	matches := getAllExercises()
	texts := map[string]string{}
	for _, value := range matches {
		index := value[1]
		texts[index] += value[0]
	}
	return texts
}

func GetAllSentences() map[string]int {
	text := GetAllTexts()
	counts := map[string]int{}
	for key, content := range text {
		rus, _ := ExtractContentSentences(content)
		counts[key] = len(rus)
	}
	return counts
}

func getAllSentencedByIndex() []string {
	byIndex := []string{}
	utils.ScanFile("data/exercise-titles.txt", func(line string) {
		line = utils.PregReplace(line, `^(\d+)`, "")
		byIndex = append(byIndex, strings.TrimSpace(line))
	})
	return byIndex
}

func GetDataForList(database *sqlx.DB) []map[string]any {
	data := make([]map[string]any, 0)
	counts := GetAllSentences()
	lines := getAllSentencedByIndex()
	for key, line := range lines {
		index := key + 1
		indexString := strconv.FormatInt(int64(index), 10)
		output := map[string]any{}
		stat := GetExerciseStatStyle(database, index)
		if stat.Page > 0 {
			a := "недавно"
			if stat.DaysAgo > 0 {
				if stat.DaysAgo > 0 {
					a = "дней:" + strconv.FormatInt(int64(stat.DaysAgo), 10)
				}
			}
			output["dateAdded"] = stat.DateAdded
			output["dateText"] = a
		}
		if counts[indexString] > 0 {
			output["timePerQuestion"] = utils.RoundFloat(stat.Time/float64(counts[indexString]), 1)
			output["errorsPerQuestion"] = utils.RoundFloat(float64(stat.Errors)/float64(counts[indexString]), 1)
		} else {
			log.Fatal("!count", indexString)
		}
		output["style"] = stat.Style
		output["time"] = stat.Time
		output["index"] = indexString
		output["count"] = counts[indexString]
		output["line"] = line
		data = append(data, output)
	}
	return data
}
