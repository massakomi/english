package utils

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func IsNumeric(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func PregReplace(content string, pattern string, replacement string) string {
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(content, replacement)
}

func PregMatchAll(pattern string, text string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindAllString(text, -1)
}

// с субпаттернами
func PregMatchAllEx(pattern string, text string) [][]string {
	re := regexp.MustCompile(pattern)
	return re.FindAllStringSubmatch(text, -1)
}

func Match(pattern string, text string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

func PregMatch(pattern string, text string) string {
	re := regexp.MustCompile(pattern)
	return re.FindString(text)
}

func PregMatchEx(pattern string, text string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindStringSubmatch(text)
}

func PregSplit(pattern string, text string) []string {
	return regexp.MustCompile(pattern).Split(text, -1)
}

func GetCookie(name string, c *gin.Context) string {
	cookie, err := c.Cookie(name)
	if err == nil {
		return ""
	}
	return cookie
}

func GetPostDefault(name string, context *gin.Context) string {
	return context.DefaultQuery(name, context.DefaultPostForm(name, GetCookie(name, context)))
}

func GetPostDefaultInt(name string, context *gin.Context) int {
	str := GetPostDefault(name, context)
	result, err := strconv.Atoi(str)
	if err != nil {
		return result
	}
	return 0
}

// MapKeySortByValues сортировка ключей словаря
func MapKeySortByValues[T int | int64](stat map[string]T, desc bool) []string {
	keys := make([]string, 0, len(stat))
	for key := range stat {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		if desc {
			return stat[keys[i]] > stat[keys[j]]
		} else {
			return stat[keys[i]] < stat[keys[j]]
		}
	})
	return keys
}

func ScanFile(filename string, callback func(string, int)) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		callback(scanner.Text(), index)
		index++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func SumIntSlice(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
