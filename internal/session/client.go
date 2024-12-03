package session

// ClientSession хранит информацию о сессии клиента.
type ClientSession struct {
	Token string `json:"token"`
}

// NewClientSession создает новый экземпляр ClientSession.
func NewClientSession() *ClientSession {
	return &ClientSession{}
}

// SetToken устанавливает токен аутентификации.
//
// token: Токен аутентификации.
func (s *ClientSession) SetToken(token string) {
	s.Token = token
}

// IsAuth проверяет, аутентифицирован ли пользователь.
//
// Возвращает true, если пользователь аутентифицирован (токен не пустой), иначе false.
func (s *ClientSession) IsAuth() bool {
	return s.Token != ""
}
