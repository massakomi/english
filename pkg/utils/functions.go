package utils

import (
	"regexp"
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
