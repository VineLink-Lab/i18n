package parser

import "golang.org/x/text/language"

const DefaultBundleName = "."

type BundleContent map[language.Tag]Contents

func NewBundleContent() BundleContent {
	return make(BundleContent)
}

func (bc *BundleContent) AddContent(tag language.Tag, content Contents) {
	(*bc)[tag] = content
}

type Contents map[string]Message

func NewContents() Contents {
	return make(Contents)
}

func (c *Contents) AddMessage(key string, message Message) {
	(*c)[key] = message
}
