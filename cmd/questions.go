package cmd

import (
	"english/pkg/utils"
	"fmt"
	_ "fmt"
	"github.com/gobs/pretty"
	"github.com/jmoiron/sqlx"
	"html/template"
	"log"
	"strconv"
	"strings"
)

func GetDataForExercise(database *sqlx.DB, exerciseIndex string) ([]map[string]any, string) {
	content := exercisesContent(exerciseIndex)
	exerciseComment := getExerciseComment(content)
	output := getOutputFromContent(content)
	outputData := make([]map[string]any, 0)
	if len(output) > 0 {
		for key, item := range output {
			exerciseStat := GetExerciseQuestionLast(database, exerciseIndex)
			outputData = append(outputData, map[string]any{
				"index":   key + 1,
				"url":     fmt.Sprintf(`http://msc/index.php?db=tester&table=english_exercise_questions&s=tbl_data&where=exercise=%v AND question=%`, exerciseIndex, key+1),
				"russian": item["rus"],
				"eng":     clearEng(item["eng"]),
				"errors":  exerciseStat[key+1]["errors"],
				"add":     template.HTML(exerciseQuestionStat(database, exerciseIndex, key+1)),
			})
		}
	}
	return outputData, exerciseComment
}

func clearEng(eng string) string {
	eng = strings.Replace(eng, "е", "e", -1)
	eng = strings.Replace(eng, `Е`, `E`, -1)
	a := utils.PregMatch(`(?i)[а-я]`, eng)
	if len(a) > 0 {
		pretty.PrettyPrint(a)
		log.Fatalf(`В eng найдены русские символы (%v)`, eng)
	}
	eng = strings.Replace(eng, `і`, `i`, -1)
	eng = strings.Replace(eng, `І`, `I`, -1)
	eng = strings.Replace(eng, `І`, `I`, -1)
	if utils.Match(`(?i)і`, eng) {
		log.Fatalf(`В eng найден i (%v)`, eng)
	}
	eng = strings.Replace(eng, `~`, `-`, -1)
	eng = strings.TrimSpace(eng)
	eng = utils.PregReplace(eng, `(?i)[\.,?!\s]$`, "")
	return eng
}

func getOutputFromContent(content string) []map[string]string {
	output := []map[string]string{}
	if content == "" {
		return output
	}
	rus, eng := ExtractContentSentences(content)
	if len(rus) != len(eng) {
		log.Fatalf(`GetDataForExercise %v != %v`, len(rus), len(eng))
	}
	for key, rusValue := range rus {
		output = append(output, map[string]string{
			"rus": rusValue,
			"eng": eng[key],
		})
	}
	return output
}

func exercisesContent(exerciseIndex string) string {
	a := getAllExercises()
	content := ""
	for _, values := range a {
		if values[1] == exerciseIndex {
			content += values[0]
		}
	}
	return content
}

func getExerciseComment(content string) string {
	exerciseComment := ""
	matches := utils.PregMatchEx(`(?i)Еxеrcіsе (\d+)`, content)
	if matches != nil {
		sents := getAllSentencedByIndex()
		i, _ := strconv.Atoi(matches[1])
		exerciseComment = sents[i-1]
	}
	return exerciseComment
}
