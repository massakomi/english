package utils

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"sort"
	"strconv"
)

func IsNumeric(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}

func PregReplace(content string, pattern string, replacement string) string {
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(content, replacement)
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
