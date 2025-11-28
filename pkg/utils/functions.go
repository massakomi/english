package utils

import (
	"github.com/gin-gonic/gin"
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

func GetCookie(name string, c *gin.Context) string {
	cookie, err := c.Cookie(name)
	if err == nil {
		return ""
	}
	return cookie
}
