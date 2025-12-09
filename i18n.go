package i18n

import (
	"github.com/VineLink-Lab/i18n/internal/translator"
)

type Translator = translator.Translator
type Option = translator.Option
type M = translator.M

var NewTranslator = translator.NewTranslator
var NewTranslatorFromFS = translator.NewTranslatorFromFS
var WithLanguage = translator.WithLanguage
var WithBundle = translator.WithBundle
var WithParams = translator.WithParams
var WithCount = translator.WithCount
