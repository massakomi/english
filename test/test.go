package test

import (
	"english/pkg/text"
	"fmt"
)

func TestGo() {

	base := text.Fs(200)
	fmt.Println(base)

	base = text.Fs(200, true, true)
	fmt.Println(base)

}

func TestBaseForm() {
	word := "upcoming"

	base := text.BaseForm(word)
	fmt.Println(base)

	base = text.BaseForm(word, true, true)
	fmt.Println(base)
}
