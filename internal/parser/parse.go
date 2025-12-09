package parser

import (
	"errors"
	"io/fs"
	"os"

	"golang.org/x/text/language"
)

func (p *Parser) PreParse() error {
	p.locker.Lock()
	defer p.locker.Unlock()
	fileInfo, err := fs.Stat(p.directory, ".")
	if err != nil {
		return err
	}
	p.bundleDir = make(map[string]string)
	if fileInfo.IsDir() {
		dirEntries, err := fs.ReadDir(p.directory, ".")
		if err != nil {
			return err
		}
		p.bundleDir = make(map[string]string)
		for _, entry := range dirEntries {
			if entry.IsDir() {
				bundlePath := entry.Name()
				p.bundleDir[entry.Name()] = bundlePath

				bundleDirEntries, err := fs.ReadDir(p.directory, bundlePath)
				if err != nil {
					return err
				}
				for _, bundleDirEntry := range bundleDirEntries {
					languageTag := parseFileNameLangTag(bundleDirEntry.Name())
					p.languages.Add(languageTag)
				}
			} else {
				p.bundleDir[DefaultBundleName] = "."
				languageTag := parseFileNameLangTag(entry.Name())
				p.languages.Add(languageTag)
			}
		}
	}
	return nil
}

func (p *Parser) ParseAllBundles() error {
	for bundleName := range p.bundleDir {
		err := p.ParseBundle(bundleName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) ParseBundle(bundleName string) error {
	for langTag := range p.languages {
		err := p.ParseContent(bundleName, langTag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) ParseContent(bundleName string, languageTag language.Tag) error {
	p.locker.Lock()
	defer p.locker.Unlock()

	bundlePath, exists := p.bundleDir[bundleName]
	if !exists {
		return errors.New("bundle not found: " + bundleName)
	}

	bundleEntries, err := fs.ReadDir(p.directory, bundlePath)
	if err != nil {
		return err
	}
	for _, entry := range bundleEntries {
		if entry.IsDir() {
			continue
		}
		if languageTag == parseFileNameLangTag(entry.Name()) {
			basePath := bundlePath + string(os.PathSeparator)
			if bundlePath == DefaultBundleName {
				basePath = ""
			}
			filePath := basePath + entry.Name()
			fileBytes, err := fs.ReadFile(p.directory, filePath)
			if err != nil {
				return err
			}
			content, err := parseFileContent(fileBytes)
			if err != nil {
				return err
			}
			bundleContent, exists := p.contents[bundleName]
			if !exists {
				bundleContent = NewBundleContent()
			}
			bundleContent.AddContent(languageTag, content)
			p.contents[bundleName] = bundleContent
		}
	}
	return nil
}

func (p *Parser) MatchLanguage(lang language.Tag) language.Tag {
	supportedLanguages := p.GetSupportedLanguages()
	matcher := language.NewMatcher(supportedLanguages)
	_, index, _ := matcher.Match(lang)
	return supportedLanguages[index]
}
