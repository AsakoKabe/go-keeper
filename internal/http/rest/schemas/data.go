package schemas

type DataResponse struct {
	Data interface{} `json:"data"`
	Meta string      `json:"meta"`
	Type string      `json:"type"`
	ID   string      `json:"ID"`
}
