package schemas

import "go-keeper/internal/models"

type FileSchema struct {
	Data models.FileData `json:"data"`
	Meta string          `json:"meta"`
}

func (f *FileSchema) Valid() bool {
	if len(f.Data.Content) == 0 || f.Data.Name == "" {
		return false
	}

	return true
}

func (f *FileSchema) GetMeta() string {
	return f.Meta
}

func (f *FileSchema) GetData() interface{} {
	return f.Data
}
