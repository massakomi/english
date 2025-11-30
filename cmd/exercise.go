package cmd

import (
	"english/pkg/utils"
	"fmt"
	"log"
	"os"
	"strings"
)

func getAllExercises() [][]string {
	content, err := os.ReadFile(`data/exercise.txt`)
	if err != nil {
		log.Fatal(err)
	}
	data := utils.PregMatchAllEx(`(?is)Еxеrcіsе (\d+)[\s\.]+[^\r\n]+`, string(content))
	return data
}

func extractSentences(rus string) []string {
	rus = strings.TrimSpace(rus)
	a := utils.PregMatchAllEx(`(?i)(\d+)[,\.]\s*(.*?)(?=\s\d{1,2}[,\.][^a-z]|$)`, rus)
	fmt.Println(a)
	data := []string{}
	return data
}

func extractContentSentences(content string) ([]string, []string) {
	content = strings.TrimSpace(content)
	if utils.PregMatch(`(i)Еxеrcіsе (\d+)`, content) != "" {
		content = utils.PregReplace(content, `(i)Еxеrcіsе (\d+)[\s\.]+`, "")
	}
	content = strings.Replace(content, "\r", "", -1)
	content = utils.PregReplace(content, `(i)([\d])t`, "$1.")
	content = utils.PregReplace(content, `(s)\n[ \t]+`, "\n")
	data := utils.PregSplit(`(s)[\r\n]{2,}`, content)
	rus := extractSentences(data[0])
	eng := extractSentences(data[1])
	return rus, eng
}

func getAllTexts() map[string]string {
	matches := getAllExercises()
	texts := map[string]string{}
	for _, value := range matches {
		index := value[1]
		texts[index] += value[0]
	}
	return texts
}

func getAllSentences() map[string]int {
	text := getAllTexts()
	counts := map[string]int{}
	for key, content := range text {
		rus, _ := extractContentSentences(content)
		counts[key] = len(rus)
	}
	return counts
}
