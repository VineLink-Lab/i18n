package main

import (
	"embed"
	"fmt"

	"github.com/VineLink-Lab/i18n/pkg/i18n"
	"golang.org/x/text/language"
)

//go:embed message
var messageFS embed.FS

func main() {
	// use with directory path
	i, err := i18n.NewI18n("example/message", language.English)
	if err != nil {
		panic(err)
	}
	testTranslation(i)

	// use with embed.FS
	i, err = i18n.NewI18nFromFS(messageFS, language.English)
	if err != nil {
		panic(err)
	}
	testTranslation(i)
}

func testTranslation(i *i18n.I18n) {
	// simple translation
	out := i.Translate("hello")
	fmt.Printf("Translation: %s\n", out)

	// translation with different language
	out = i.Translate("hello", i18n.WithLanguage(language.French))
	fmt.Printf("Translation: %s\n", out)

	// translation with non-supported language falls back to default
	out = i.Translate("hello", i18n.WithLanguage(language.Spanish))
	fmt.Printf("Translation: %s\n", out)

	// translation with child language falls back to parent language
	out = i.Translate("hello", i18n.WithLanguage(language.AmericanEnglish))
	fmt.Printf("Translation: %s\n", out)

	// choosing different bundle
	i.SetDefaultBundle("with_parameters")

	// translation with parameters
	out = i.Translate("hello", i18n.WithParams(i18n.M{"name": "Alice"}))
	fmt.Printf("Translation: %s\n", out)

	// translation with parameters and different language
	out = i.Translate("hello", i18n.WithLanguage(language.French), i18n.WithParams(i18n.M{"name": "Alice"}))
	fmt.Printf("Translation: %s\n", out)
}
