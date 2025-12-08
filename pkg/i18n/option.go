package i18n

import "golang.org/x/text/language"

type TranslationOptions struct {
	Lang   language.Tag
	Bundle string
	Params map[string]any
	Count  int
}
type Option func(options *TranslationOptions)

type M map[string]any

func WithLanguage(lang language.Tag) Option {
	return func(i *TranslationOptions) {
		i.Lang = lang
	}
}

func WithBundle(bundle string) Option {
	return func(i *TranslationOptions) {
		i.Bundle = bundle
	}
}

func WithParams(params map[string]any) Option {
	return func(i *TranslationOptions) {
		i.Params = params
	}
}

func WithCount(count int) Option {
	return func(i *TranslationOptions) {
		i.Count = count
	}
}
