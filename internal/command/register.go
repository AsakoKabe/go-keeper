package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go-keeper/internal/models"
	"go-keeper/internal/session"
	"go-keeper/pkg/middleware"
)

// RegisterCMD регистрация пользователя
type RegisterCMD struct {
	httpClient    *http.Client
	clientSession *session.ClientSession
	serverAddr    string
}

// NewRegisterCMD создает новый экземпляр RegisterCMD.
//
// httpClient: HTTP-клиент для выполнения запросов.
// clientSession: Сессия клиента для хранения токена.
// serverAddr: Адрес сервера.
func NewRegisterCMD(
	httpClient *http.Client, clientSession *session.ClientSession, serverAddr string,
) *RegisterCMD {
	return &RegisterCMD{
		httpClient:    httpClient,
		clientSession: clientSession,
		serverAddr:    serverAddr,
	}
}

// Execute выполняет команду регистрации.
//
// args: Аргументы команды, где args[0] - логин, args[1] - пароль.
//
// Возвращает ошибку, если регистрация не удалась.
// В случае успешной регистрации токен сохраняется в сессии клиента.
func (r *RegisterCMD) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("failed to register, not enough arguments")
	}

	body := models.User{
		Login:    args[0],
		Password: args[1],
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := r.httpClient.Post(
		fmt.Sprintf("%s/user/register", r.serverAddr),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register, %d", resp.StatusCode)
	}

	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == middleware.CookieName {
			r.clientSession.SetToken(c.Value)
		}
	}

	return nil
}
