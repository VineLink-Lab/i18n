package parser

import (
	"errors"
	"io/fs"
	"os"
	"sync"

	"github.com/VineLink-Lab/i18n/utils"

	"golang.org/x/text/language"
)

type Parser struct {
	directoryPath string
	directory     fs.FS
	bundleDir     map[string]string
	languages     utils.Set[language.Tag]

	contents map[string]BundleContent

	locker sync.Mutex
}

func NewParser(directory string) (*Parser, error) {
	dir := os.DirFS(directory)
	return NewParserFromFS(dir)
}

func NewParserFromFS(directory fs.FS) (*Parser, error) {
	p := &Parser{
		directory: directory,
		bundleDir: make(map[string]string),
		languages: utils.NewSet[language.Tag](),
		contents:  make(map[string]BundleContent),
	}
	err := p.PreParse()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) GetSupportedLanguages() []language.Tag {
	return p.languages.ToSlice()
}

func (p *Parser) Translate(key string, lang language.Tag, bundle string, count int, params map[string]any) (string, error) {
	if bundle == "" {
		bundle = DefaultBundleName
	}
	lang = p.MatchLanguage(lang)
	var err error
	contents, ok := p.contents[bundle]
	if !ok {
		err = p.ParseContent(bundle, lang)
		if err != nil {
			return "", err
		}
		contents, ok = p.contents[bundle]
		if !ok {
			return "", errors.New("bundle not found")
		}
	}
	langContents, ok := contents[lang]
	if !ok {
		err = p.ParseContent(bundle, lang)
		if err != nil {
			return "", err
		}
		langContents, ok = contents[lang]
		if !ok {
			return "", errors.New("language not found in bundle")
		}
	}
	message, ok := langContents[key]
	if !ok {
		return "", errors.New("message key not found")
	}
	return message.Translate(count, params, lang), nil
}

func (p *Parser) Write() error {
	panic("not implemented")
}
