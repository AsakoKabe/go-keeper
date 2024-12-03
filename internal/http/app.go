package http

import (
	"go-keeper/config/server"
)

// App интерфейс для работы с приложением
type App interface {
	Run(cfg *server.Config) error
	Stop()
}
