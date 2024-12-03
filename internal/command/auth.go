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

// AuthCMD команда для аутентификации пользователя.
type AuthCMD struct {
	httpClient    *http.Client
	clientSession *session.ClientSession
	serverAddr    string
}

// NewAuthCMD создает новый экземпляр AuthCMD.
//
// httpClient: HTTP-клиент для выполнения запросов.
// clientSession: Сессия клиента для хранения токена.
// serverAddr: Адрес сервера.
func NewAuthCMD(
	httpClient *http.Client, clientSession *session.ClientSession, serverAddr string,
) *AuthCMD {
	return &AuthCMD{
		httpClient:    httpClient,
		clientSession: clientSession,
		serverAddr:    serverAddr,
	}
}

// Execute выполняет команду аутентификации.
//
// args: Аргументы команды, где args[0] - логин, args[1] - пароль.
//
// Возвращает ошибку, если аутентификация не удалась.
// В случае успешной аутентификации токен сохраняется в сессии клиента.
func (c *AuthCMD) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("failed to auth, not enough arguments")
	}

	body := models.User{
		Login:    args[0],
		Password: args[1],
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/user/auth", c.serverAddr),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to auth, %d", resp.StatusCode)
	}

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == middleware.CookieName {
			c.clientSession.SetToken(cookie.Value)
		}
	}

	return nil
}
