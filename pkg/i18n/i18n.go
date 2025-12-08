package i18n

import (
	"io/fs"

	"github.com/VineLink-Lab/i18n/parser"
	"golang.org/x/text/language"
)

type I18n struct {
	parser        *parser.Parser
	defaultLang   language.Tag
	defaultBundle string
}

func NewI18n(directory string, defaultLanguage ...language.Tag) (*I18n, error) {
	p, err := parser.NewParser(directory)
	if err != nil {
		return nil, err
	}
	return newI18n(p, defaultLanguage...), nil
}

func NewI18nFromFS(directory fs.FS, defaultLanguage ...language.Tag) (*I18n, error) {
	rootFS, err := fs.ReadDir(directory, ".")
	if err != nil {
		return nil, err
	}
	sbuFS, err := fs.Sub(directory, rootFS[0].Name())
	if err != nil {
		return nil, err
	}
	p, err := parser.NewParserFromFS(sbuFS)
	if err != nil {
		return nil, err
	}
	return newI18n(p, defaultLanguage...), nil
}

func newI18n(p *parser.Parser, defaultLanguage ...language.Tag) *I18n {
	supportedLanguages := p.GetSupportedLanguages()
	if len(supportedLanguages) == 0 {
		supportedLanguages = []language.Tag{language.English}
	}
	i := &I18n{
		parser:        p,
		defaultLang:   supportedLanguages[0],
		defaultBundle: parser.DefaultBundleName,
	}
	if len(defaultLanguage) > 0 {
		i.SetDefaultLanguage(defaultLanguage[0])
	}
	return i
}

func (i *I18n) SetDefaultLanguage(lang language.Tag) {
	i.defaultLang = lang
}

func (i *I18n) GetDefaultLanguage() language.Tag {
	return i.defaultLang
}

func (i *I18n) SetDefaultBundle(bundle string) {
	i.defaultBundle = bundle
}

func (i *I18n) GetDefaultBundle() string {
	return i.defaultBundle
}

func (i *I18n) Clone() *I18n {
	clone := *i
	return &clone
}

func (i *I18n) GetSupportedLanguages() []language.Tag {
	return i.parser.GetSupportedLanguages()
}

func (i *I18n) MatchLanguage(languageCode string) language.Tag {
	lang, err := language.Parse(languageCode)
	if err != nil {
		return language.Und
	}
	return i.parser.MatchLanguage(lang)
}

func (i *I18n) TranslateE(key string, options ...Option) (string, error) {
	var opts TranslationOptions
	for _, opt := range options {
		opt(&opts)
	}
	return i.TranslateWithOptions(key, opts)
}

func (i *I18n) Translate(key string, options ...Option) string {
	out, err := i.TranslateE(key, options...)
	if err != nil {
		return key
	}
	return out
}

func (i *I18n) TranslateWithOptions(key string, options TranslationOptions) (string, error) {
	if options.Lang == language.Und {
		options.Lang = i.defaultLang
	}
	if options.Bundle == "" {
		options.Bundle = i.defaultBundle
	}
	return i.parser.Translate(key, options.Lang, options.Bundle, options.Count, options.Params)
}
