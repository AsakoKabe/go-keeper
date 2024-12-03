package models

type Data struct {
	Data string `json:"data"`
	Meta string `json:"meta"`
	Type string `json:"type"`
	ID   string `json:"ID"`
}

type LogPassData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CardData struct {
	Number    string `json:"number"`
	ExpiredAt string `json:"expired_at"`
	CVV       string `json:"cvv"`
}

type TextData struct {
	Text string `json:"text"`
}

type FileData struct {
	Content []byte `json:"content"`
	Name    string `json:"name"`
}
