package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	contextUtils "go-keeper/internal/context"
	"go-keeper/internal/db/repository/errs"
	"go-keeper/internal/http/rest/schemas"
	"go-keeper/internal/keeper/data"
	"go-keeper/internal/utils"
)

// Data интерфейс работы с данными
type Data interface {
	GetData() interface{}
	GetMeta() string
	Valid() bool
}

// DataHandler Структура для endpoints
type DataHandler struct {
	dataService *data.Service
}

// NewDataHandler создает новый экземпляр DataHandler.
func NewDataHandler(dataService *data.Service) *DataHandler {
	return &DataHandler{dataService: dataService}
}

// Add обрабатывает POST-запрос на добавление данных.
func (h *DataHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID := contextUtils.GetUserID(r.Context())

	dataType := chi.URLParam(r, "type")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error to read body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataBody, err := h.parseDataRequest(body, dataType)
	if err != nil {
		slog.Error("error to parse body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !dataBody.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := h.dataService.Add(
		r.Context(), userID, dataType, dataBody.GetData(), dataBody.GetMeta(),
	)
	if err != nil {
		slog.Error("error to save data", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	dataResponse := schemas.DataResponse{
		Data: dataBody.GetData(),
		Meta: dataBody.GetMeta(),
		ID:   id,
		Type: dataType,
	}
	err = json.NewEncoder(w).Encode(dataResponse)
	if err != nil {
		slog.Error("error to serialize response", slog.String("err", err.Error()))
		return
	}
}

// GetAllData обрабатывает GET-запрос на получение всех данных определенного типа для пользователя
func (h *DataHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	userID := contextUtils.GetUserID(r.Context())
	dataType := chi.URLParam(r, "type")
	userDataSet, err := h.dataService.GetAllData(r.Context(), userID, dataType)
	if err != nil {
		slog.Error("error to get all data", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userDataSet)
	if err != nil {
		slog.Error("error to serialize response", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GeByID обрабатывает GET-запрос на получение данных по ID.
func (h *DataHandler) GeByID(w http.ResponseWriter, r *http.Request) {
	userID := contextUtils.GetUserID(r.Context())
	dataID := chi.URLParam(r, "dataID")
	if !utils.IsValidUUID(dataID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dataType := chi.URLParam(r, "type")

	userData, err := h.dataService.GetByID(r.Context(), userID, dataID, dataType)
	if errors.Is(err, errs.ErrDataNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("error to get data by id", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userData)
	if err != nil {
		slog.Error("error to serialize response", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// DeleteByID обрабатывает DELETE-запрос на удаление данных по ID.
func (h *DataHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	userID := contextUtils.GetUserID(r.Context())
	dataID := chi.URLParam(r, "dataID")
	if !utils.IsValidUUID(dataID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dataType := chi.URLParam(r, "type")

	err := h.dataService.DeleteByID(r.Context(), userID, dataID, dataType)
	if errors.Is(err, errs.ErrDataNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("error to delete data by id", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Update обрабатывает PUT-запрос на обновление данных.
func (h *DataHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := contextUtils.GetUserID(r.Context())
	dataID := chi.URLParam(r, "dataID")
	if !utils.IsValidUUID(dataID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dataType := chi.URLParam(r, "type")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error to read body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataBody, err := h.parseDataRequest(body, dataType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !dataBody.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.dataService.Update(
		r.Context(), userID, dataID, dataType, dataBody.GetData(), dataBody.GetMeta(),
	)
	if errors.Is(errs.ErrDataNotFound, err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("error to update data", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dataResponse := schemas.DataResponse{
		Data: dataBody.GetData(),
		Meta: dataBody.GetMeta(),
		ID:   dataID,
		Type: dataType,
	}
	err = json.NewEncoder(w).Encode(dataResponse)
	if err != nil {
		slog.Error("error to serialize response", slog.String("err", err.Error()))
		return
	}
}

// parseDataRequest разбирает тело запроса в зависимости от типа данных.
func (h *DataHandler) parseDataRequest(
	body []byte, dataType string,
) (Data, error) {

	switch dataType {
	case data.LogPassDataType:
		var logpass schemas.LogPassSchema
		if err := json.Unmarshal(body, &logpass); err != nil {
			return nil, err
		}

		return &logpass, nil
	case data.CardDataType:
		var card schemas.CardSchema
		if err := json.Unmarshal(body, &card); err != nil {
			return nil, err
		}
		return &card, nil
	case data.TextDataType:
		var text schemas.TextSchema
		if err := json.Unmarshal(body, &text); err != nil {
			return nil, err
		}
		return &text, nil
	case data.FileDataType:
		var file schemas.FileSchema
		if err := json.Unmarshal(body, &file); err != nil {
			return nil, err
		}
		return &file, nil
	}

	return nil, ErrInvalidDataBody
}
