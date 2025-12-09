package translator

import (
	"io/fs"

	"github.com/VineLink-Lab/i18n/internal/parser"
	"golang.org/x/text/language"
)

type Translator struct {
	parser        *parser.Parser
	defaultLang   language.Tag
	defaultBundle string
}

func NewTranslator(directory string, defaultLanguage ...language.Tag) (*Translator, error) {
	p, err := parser.NewParser(directory)
	if err != nil {
		return nil, err
	}
	return newTranslator(p, defaultLanguage...), nil
}

func NewTranslatorFromFS(directory fs.FS, defaultLanguage ...language.Tag) (*Translator, error) {
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
	return newTranslator(p, defaultLanguage...), nil
}

func newTranslator(p *parser.Parser, defaultLanguage ...language.Tag) *Translator {
	supportedLanguages := p.GetSupportedLanguages()
	if len(supportedLanguages) == 0 {
		supportedLanguages = []language.Tag{language.English}
	}
	t := &Translator{
		parser:        p,
		defaultLang:   supportedLanguages[0],
		defaultBundle: parser.DefaultBundleName,
	}
	if len(defaultLanguage) > 0 {
		t.defaultLang = defaultLanguage[0]
	}
	return t
}

func (t *Translator) UseLanguageCode(langCode string) *Translator {
	langTag := t.MatchLanguageCode(langCode)
	return t.UseLanguage(langTag)
}

func (t *Translator) UseLanguage(lang language.Tag) *Translator {
	tc := t.Clone()
	tc.defaultLang = lang
	return tc
}

func (t *Translator) Language() language.Tag {
	return t.defaultLang
}

func (t *Translator) UseBundle(bundle string) *Translator {
	tc := t.Clone()
	tc.defaultBundle = bundle
	return tc
}

func (t *Translator) Bundle() string {
	return t.defaultBundle
}

func (t *Translator) Clone() *Translator {
	clone := *t
	return &clone
}

func (t *Translator) GetSupportedLanguages() []language.Tag {
	return t.parser.GetSupportedLanguages()
}

func (t *Translator) MatchLanguageCode(languageCode string) language.Tag {
	lang, err := language.Parse(languageCode)
	if err != nil {
		return language.Und
	}
	return t.parser.MatchLanguage(lang)
}

func (t *Translator) MatchLanguage(languageTag language.Tag) language.Tag {
	return t.parser.MatchLanguage(languageTag)
}

func (t *Translator) TranslateE(key string, options ...Option) (string, error) {
	var opts Options
	for _, opt := range options {
		opt(&opts)
	}
	return t.TranslateWithOptions(key, opts)
}

func (t *Translator) Translate(key string, options ...Option) string {
	out, err := t.TranslateE(key, options...)
	if err != nil {
		return key
	}
	return out
}

func (t *Translator) TranslateWithOptions(key string, options Options) (string, error) {
	if options.Lang == language.Und {
		options.Lang = t.defaultLang
	}
	if options.Bundle == "" {
		options.Bundle = t.defaultBundle
	}
	return t.parser.Translate(key, options.Lang, options.Bundle, options.Count, options.Params)
}
