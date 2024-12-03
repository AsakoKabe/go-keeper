package schemas

import "go-keeper/internal/models"

type LogPassSchema struct {
	Data models.LogPassData `json:"data"`
	Meta string             `json:"meta"`
}

func (l *LogPassSchema) Valid() bool {
	if l.Data.Login == "" || l.Data.Password == "" {
		return false
	}

	return true
}

func (l *LogPassSchema) GetData() interface{} {
	return l.Data
}

func (l *LogPassSchema) GetMeta() string {
	return l.Meta
}
