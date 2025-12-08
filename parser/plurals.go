package parser

import "golang.org/x/text/language"

type Plurals struct {
	Zero  string `json:"zero"`
	One   string `json:"one"`
	Two   string `json:"two"`
	Few   string `json:"few"`
	Many  string `json:"many"`
	Other string `json:"other"`
}

func (p Plurals) GetContent(languageTag language.Tag, count int) string {
	return p.Other
}
