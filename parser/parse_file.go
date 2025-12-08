package parser

import (
	"encoding/json"
	"log/slog"
)

type fileContent map[string]json.RawMessage

type contentItem struct {
	Content     string   `json:"content"`
	Description string   `json:"description"`
	Plurals     *Plurals `json:"plurals"`
}

func parseFileContent(fileBytes []byte) (Contents, error) {
	content := NewContents()
	var fc fileContent
	err := json.Unmarshal(fileBytes, &fc)
	if err != nil {
		return nil, err
	}
	for key, rawMsg := range fc {
		var message Message

		var value string
		err := json.Unmarshal(rawMsg, &value)
		if err == nil {
			message = Message{
				Content: value,
			}
			content.AddMessage(key, message)
		} else {
			var item contentItem
			err := json.Unmarshal(rawMsg, &item)
			if err != nil {
				slog.Error("failed to unmarshal content item", slog.String("key", key), slog.String("error", err.Error()))
				continue
			}
			message = Message{
				Content:     item.Content,
				Description: item.Description,
				Plurals:     item.Plurals,
			}
			content.AddMessage(key, message)
		}
	}

	return content, nil
}
