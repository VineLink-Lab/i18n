package translator

import "golang.org/x/text/language"

type Options struct {
	Lang   language.Tag
	Bundle string
	Params map[string]any
	Count  int
}
type Option func(options *Options)

type M map[string]any

func WithLanguage(lang language.Tag) Option {
	return func(i *Options) {
		i.Lang = lang
	}
}

func WithBundle(bundle string) Option {
	return func(i *Options) {
		i.Bundle = bundle
	}
}

func WithParams(params map[string]any) Option {
	return func(i *Options) {
		i.Params = params
	}
}

func WithCount(count int) Option {
	return func(i *Options) {
		i.Count = count
	}
}
