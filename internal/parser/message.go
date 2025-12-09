package parser

import (
	"strings"

	"github.com/spf13/cast"
	"golang.org/x/text/language"
)

type Message struct {
	Content     string   `json:"content"`
	Description string   `json:"description"`
	Plurals     *Plurals `json:"plurals"`
}

func (m Message) Translate(count int, params map[string]any, languageTag language.Tag) string {
	msg := m.getPluralForm(count, languageTag)
	// Simple parameter replacement
	for k, v := range params {
		placeholder := "{{" + k + "}}"
		msg = strings.Replace(msg, placeholder, cast.ToString(v), -1)
	}
	return msg
}

func (m Message) getPluralForm(count int, languageTag language.Tag) string {
	msg := m.Content
	if m.Plurals != nil {
		msg = m.Plurals.GetContent(languageTag, count)
	}
	return msg
}
