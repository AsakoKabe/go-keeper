package schemas

import "go-keeper/internal/models"

type TextSchema struct {
	Data models.TextData `json:"data"`
	Meta string          `json:"meta"`
}

func (t *TextSchema) Valid() bool {
	if t.Data.Text == "" {
		return false
	}

	return true
}

func (t *TextSchema) GetMeta() string {
	return t.Meta
}

func (t *TextSchema) GetData() interface{} {
	return t.Data
}
