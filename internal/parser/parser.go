package parser

import (
	"encoding/json"
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
	p, err := NewParserFromFS(dir)
	if err != nil {
		return nil, err
	}
	p.directoryPath = directory
	return p, nil
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

func (p *Parser) GetDirectoryPath() string {
	return p.directoryPath
}

func (p *Parser) GetAvailableBundles() []string {
	bundles := make([]string, 0, len(p.bundleDir))
	for bundle := range p.bundleDir {
		bundles = append(bundles, bundle)
	}
	return bundles
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

func (p *Parser) GetBundleContent(bundle string) (BundleContent, error) {
	if bundle == "" {
		bundle = DefaultBundleName
	}

	// 如果内容已经加载,直接返回
	if content, ok := p.contents[bundle]; ok {
		return content, nil
	}

	// 加载所有语言的内容
	for _, lang := range p.GetSupportedLanguages() {
		err := p.ParseContent(bundle, lang)
		if err != nil {
			return nil, err
		}
	}

	content, ok := p.contents[bundle]
	if !ok {
		return nil, errors.New("bundle not found")
	}

	return content, nil
}

// UpdateMessage updates or adds a translation key
func (p *Parser) UpdateMessage(bundle string, key string, lang language.Tag, message Message) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	if bundle == "" {
		bundle = DefaultBundleName
	}

	// Ensure bundle content is loaded
	bundleContent, ok := p.contents[bundle]
	if !ok {
		bundleContent = NewBundleContent()
		p.contents[bundle] = bundleContent
	}

	// Ensure language content is loaded
	langContent, ok := bundleContent[lang]
	if !ok {
		langContent = NewContents()
		bundleContent[lang] = langContent
	}

	// Update message
	langContent[key] = message

	// Write to file immediately
	return p.writeLanguageFile(bundle, lang)
}

// DeleteMessage deletes a translation key
func (p *Parser) DeleteMessage(bundle string, key string) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	if bundle == "" {
		bundle = DefaultBundleName
	}

	bundleContent, ok := p.contents[bundle]
	if !ok {
		return errors.New("bundle not found")
	}

	// Delete the key from all languages
	for lang, langContent := range bundleContent {
		delete(langContent, key)
		// Write each language file
		err := p.writeLanguageFile(bundle, lang)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeLanguageFile writes the content of a specific language to file
func (p *Parser) writeLanguageFile(bundle string, lang language.Tag) error {
	if p.directoryPath == "" {
		return errors.New("cannot write to embedded filesystem")
	}

	bundleContent, ok := p.contents[bundle]
	if !ok {
		return errors.New("bundle not found")
	}

	langContent, ok := bundleContent[lang]
	if !ok {
		return errors.New("language content not found")
	}

	// Build file path
	bundlePath := p.bundleDir[bundle]
	if bundlePath == "" {
		bundlePath = DefaultBundleName
	}

	var filePath string
	langStr := lang.String()
	if bundlePath == DefaultBundleName {
		filePath = p.directoryPath + string(os.PathSeparator) + langStr + ".json"
	} else {
		filePath = p.directoryPath + string(os.PathSeparator) + bundlePath + string(os.PathSeparator) + langStr + ".json"
	}

	// Convert to JSON format
	fileContent := make(map[string]interface{})
	for key, msg := range langContent {
		// If only content exists, store as string directly
		if msg.Description == "" && msg.Plurals == nil {
			fileContent[key] = msg.Content
		} else {
			// Otherwise store complete object
			item := map[string]interface{}{
				"content": msg.Content,
			}
			if msg.Description != "" {
				item["description"] = msg.Description
			}
			if msg.Plurals != nil {
				item["plurals"] = msg.Plurals
			}
			fileContent[key] = item
		}
	}

	// Serialize to JSON
	jsonBytes, err := json.MarshalIndent(fileContent, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, jsonBytes, 0644)
}

func (p *Parser) Write() error {
	panic("not implemented")
}
