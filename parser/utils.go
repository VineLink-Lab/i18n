package parser

import (
	"log/slog"
	"os"

	"golang.org/x/text/language"
)

func parseFileNameLangTag(fileName string) (langTag language.Tag) {
	var tag string
	formatStartIdx := -1
	for i := len(fileName) - 1; i >= 0; i-- {
		c := fileName[i]
		if os.IsPathSeparator(c) {
			if formatStartIdx != -1 {
				tag = fileName[i+1 : formatStartIdx]
			}
			break
		}
		if fileName[i] == '.' {
			if formatStartIdx != -1 {
				tag = fileName[i+1 : formatStartIdx]
				break
			} else {
				formatStartIdx = i
			}
		}
	}
	if formatStartIdx != -1 {
		tag = fileName[:formatStartIdx]
	}
	langTag, err := language.Parse(tag)
	if err != nil {
		slog.Error("failed to parse language tag from path", slog.String("file name", fileName), slog.String("error", err.Error()))
		return language.Und
	}
	return
}
