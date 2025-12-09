package main

import (
	"embed"
	"fmt"

	"github.com/VineLink-Lab/i18n/internal/translator"
	"golang.org/x/text/language"
)

//go:embed message
var messageFS embed.FS

func main() {
	// use with directory path
	i, err := translator.NewTranslator("example/message", language.English)
	if err != nil {
		panic(err)
	}
	testTranslation(i)

	// use with embed.FS
	i, err = translator.NewTranslatorFromFS(messageFS, language.English)
	if err != nil {
		panic(err)
	}
	testTranslation(i)
}

func testTranslation(t *translator.Translator) {
	// simple translation
	out := t.Translate("hello")
	fmt.Printf("Translation: %s\n", out)

	// translation with different language
	out = t.Translate("hello", translator.WithLanguage(language.French))
	fmt.Printf("Translation: %s\n", out)

	// translation with non-supported language falls back to default
	out = t.Translate("hello", translator.WithLanguage(language.Spanish))
	fmt.Printf("Translation: %s\n", out)

	// translation with child language falls back to parent language
	out = t.Translate("hello", translator.WithLanguage(language.AmericanEnglish))
	fmt.Printf("Translation: %s\n", out)

	// choosing different bundle
	t = t.UseBundle("with_parameters")

	// translation with parameters
	out = t.Translate("hello", translator.WithParams(translator.M{"name": "Alice"}))
	fmt.Printf("Translation: %s\n", out)

	out = t.UseLanguage(language.Spanish).Translate("hello", translator.WithParams(translator.M{"name": "Alice"}))
	fmt.Printf("Translation: %s\n", out)

	// translation with parameters and different language
	out = t.Translate("hello", translator.WithLanguage(language.French), translator.WithParams(translator.M{"name": "Alice"}))
	fmt.Printf("Translation: %s\n", out)
}
