package text

import (
	"english/pkg/utils"
	"fmt"
	_ "fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
)

// BaseForm
// space - взять только первое слово до первого пробела (например если передается фраза)
func BaseForm(es string, space ...bool) string {

	if strings.Contains(es, " ") && len(space) > 0 {
		index := strings.Index(es, " ")
		s := string([]rune(es)[0:index])
		if len(s) > 3 {
			es = s
		}
	}
	esx := utils.PregReplace(es, `(less|ness|ance|able|ment|ing|ful|ly)$`, "")
	esx = utils.PregReplace(esx, `([^a])y$`, "$1")
	esx = strings.Replace(esx, `ation`, `ate`, -1)
	if len(esx) > 3 {
		// deterioration --> deteriorate
		if esx != es {
			// jarr - jar
			last1 := esx[len(esx)-1:]
			last2 := esx[len(esx)-2 : len(esx)-1]
			if last1 == last2 {
				esx = esx[0 : len(esx)-1]
			}
			es = esx
		}
	}
	if len(es) > 5 {
		if strings.HasSuffix(esx, "ed") {
			es = es[0 : len(esx)-1]
		}
	}

	return esx
}

// Fs format seconds
func Fs(seconds float64, opts ...bool) template.HTML {
	s := false
	h := true
	if len(opts) == 1 {
		s = opts[0]
	}
	if len(opts) == 2 {
		s = opts[0]
		h = opts[1]
	}
	hs := math.Floor(seconds / 3600)
	ms := math.Floor((seconds - hs*60) / 60)
	var output string
	if h {
		output = fmt.Sprintf("%v:%02d", hs, int(ms))
	} else {
		output = strconv.FormatInt(int64(ms), 10)
	}
	if s {
		ss := math.Floor(seconds - hs*3600 - ms*60)
		output += fmt.Sprintf(":%02d", int(ss))
	}
	return template.HTML(fmt.Sprintf(`<span title="%v">%v</span>`, seconds, output))
}
