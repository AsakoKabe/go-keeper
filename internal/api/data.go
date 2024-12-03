package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go-keeper/internal/http/rest/schemas"
	"go-keeper/internal/keeper/data"
	"go-keeper/internal/session"
	"go-keeper/pkg/middleware"
)

// DataAPI предоставляет методы для взаимодействия с API данных.
type DataAPI struct {
	baseHTTPAddress string
	httpClient      *http.Client
	session         *session.ClientSession
}

// NewDataAPI создает новый экземпляр DataAPI.
//
// baseHTTPAddress: Базовый адрес HTTP сервера.
// httpClient: HTTP клиент для выполнения запросов.
// session: Сессия клиента для аутентификации.
func NewDataAPI(
	baseHTTPAddress string,
	httpClient *http.Client,
	session *session.ClientSession,
) *DataAPI {
	return &DataAPI{
		baseHTTPAddress: baseHTTPAddress,
		httpClient:      httpClient,
		session:         session,
	}
}

// GetAll получает все данные указанного типа.
//
// dataType: Тип данных для получения.
//
// Возвращает слайс структур schemas.DataResponse, содержащий полученные данные, и ошибку, если таковая возникла.
func (api *DataAPI) GetAll(dataType string) ([]schemas.DataResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf("%s/user/data/%s", api.baseHTTPAddress, dataType), nil,
	)
	if err != nil {
		return nil, err
	}

	req.AddCookie(
		&http.Cookie{
			Name:  middleware.CookieName,
			Value: api.session.Token,
		},
	)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get all data, %d", resp.StatusCode)
	}

	var respBody []schemas.DataResponse
	if err = json.Unmarshal(body, &respBody); err != nil {
		return nil, err
	}

	for idx, userData := range respBody {
		respBody[idx].Data, err = api.parseData(userData)
		if err != nil {
			return nil, err
		}
	}

	return respBody, nil
}

// GetByID получает данные по указанному ID и типу.
//
// dataID: ID данных для получения.
// dataType: Тип данных для получения.
//
// Возвращает указатель на структуру schemas.DataResponse, содержащую полученные данные, и ошибку, если таковая возникла.
func (api *DataAPI) GetByID(dataID string, dataType string) (*schemas.DataResponse, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/user/data/%s/%s", api.baseHTTPAddress, dataType, dataID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.AddCookie(
		&http.Cookie{
			Name:  middleware.CookieName,
			Value: api.session.Token,
		},
	)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data by id, %d", resp.StatusCode)
	}

	var respBody schemas.DataResponse
	if err = json.Unmarshal(body, &respBody); err != nil {
		return nil, err
	}
	respBody.Data, err = api.parseData(respBody)
	if err != nil {
		return nil, err
	}

	return &respBody, nil
}

// DeleteByID удаляет данные по указанному ID и типу.
//
// dataID: ID данных для удаления.
// dataType: Тип данных для удаления.
//
// Возвращает ошибку, если таковая возникла.
func (api *DataAPI) DeleteByID(dataID string, dataType string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/user/data/%s/%s", api.baseHTTPAddress, dataType, dataID),
		nil,
	)
	if err != nil {
		return err
	}

	req.AddCookie(
		&http.Cookie{
			Name:  middleware.CookieName,
			Value: api.session.Token,
		},
	)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete data by id, %d", resp.StatusCode)
	}

	return nil
}

// Add добавляет новые данные.
//
// entity: Структура schemas.DataResponse, содержащая данные для добавления.
//
// Возвращает указатель на структуру schemas.DataResponse, содержащую добавленные данные, и ошибку, если таковая возникла.
func (api *DataAPI) Add(entity schemas.DataResponse) (*schemas.DataResponse, error) {
	b, _ := json.Marshal(entity)

	req, err := http.NewRequest(
		http.MethodPost, fmt.Sprintf("%s/user/data/%s", api.baseHTTPAddress, entity.Type),
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(
		&http.Cookie{
			Name:  middleware.CookieName,
			Value: api.session.Token,
		},
	)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to add data, %d", resp.StatusCode)
	}

	var respBody schemas.DataResponse
	if err = json.Unmarshal(userData, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

// Update обновляет существующие данные.
//
// entity: Структура schemas.DataResponse, содержащая данные для обновления.
//
// Возвращает указатель на структуру schemas.DataResponse, содержащую обновленные данные, и ошибку, если таковая возникла.
func (api *DataAPI) Update(entity schemas.DataResponse) (*schemas.DataResponse, error) {
	b, _ := json.Marshal(entity)

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/user/data/%s/%s", api.baseHTTPAddress, entity.Type, entity.ID),
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(
		&http.Cookie{
			Name:  middleware.CookieName,
			Value: api.session.Token,
		},
	)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update data, %d", resp.StatusCode)
	}

	var respBody schemas.DataResponse
	if err = json.Unmarshal(userData, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

// parseData разбирает данные из JSON и преобразует их в соответствующий тип данных.
//
// userData: Структура schemas.DataResponse, содержащая данные для разбора.
//
// Возвращает интерфейс, содержащий разобранные данные, и ошибку, если таковая возникла.
func (api *DataAPI) parseData(userData schemas.DataResponse) (interface{}, error) {
	jsonData, err := json.Marshal(userData.Data)
	if err != nil {
		return nil, err
	}

	return data.ParseJsonData(jsonData, userData.Type)
}
