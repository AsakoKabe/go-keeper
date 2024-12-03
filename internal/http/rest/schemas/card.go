package schemas

import "go-keeper/internal/models"

type CardSchema struct {
	Data models.CardData `json:"data"`
	Meta string          `json:"meta"`
}

func (c *CardSchema) Valid() bool {
	if c.Data.Number == "" || c.Data.ExpiredAt == "" || c.Data.CVV == "" {
		return false
	}

	return true
}

func (c *CardSchema) GetMeta() string {
	return c.Meta
}

func (c *CardSchema) GetData() interface{} {
	return c.Data
}
