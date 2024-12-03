package command

import (
	"fmt"

	"go-keeper/internal/api"
	"go-keeper/internal/http/rest/schemas"
	"go-keeper/internal/keeper/data"
	"go-keeper/internal/models"
	"go-keeper/pkg/console"
)

const (
	TextCreate = "create"
	TextGet    = "get"
	TextGetAll = "get-all"
	TextUpdate = "update"
	TextDelete = "delete"
)

// TextCMD управления текстовыми данными
type TextCMD struct {
	subCommandStrategies map[string]subCommandStrategy
	dataAPI              *api.DataAPI
}

// NewTextCMD создает новый экземпляр TextCMD.
//
// dataAPI: API для взаимодействия с данными.
func NewTextCMD(dataAPI *api.DataAPI) *TextCMD {
	textCMD := &TextCMD{dataAPI: dataAPI}
	textCMD.subCommandStrategies = map[string]subCommandStrategy{
		TextCreate: textCMD.create,
		TextGet:    textCMD.get,
		TextUpdate: textCMD.update,
		TextDelete: textCMD.delete,
		TextGetAll: textCMD.getAll,
	}
	return textCMD
}

// Execute выполняет подкоманду для управления текстовыми данными.
//
// args: Аргументы команды, где args[0] - подкоманда (create, get, update, delete, get-all).
//
// Возвращает ошибку, если подкоманда не найдена или произошла ошибка при ее выполнении.
func (t *TextCMD) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("sub command for text not found")
	}

	subCMD, ok := t.subCommandStrategies[args[0]]
	if !ok {
		return fmt.Errorf("sub command not found")
	}

	return subCMD(args[1:])
}

// create создает новую текстовую запись.
//
// _: Аргументы подкоманды (не используются).
//
// Запрашивает у пользователя текст и метаданные, и сохраняет их через API.
// Возвращает ошибку, если не удалось получить данные или сохранить запись.
func (t *TextCMD) create(_ []string) error {
	text := console.GetInput("Enter text: ", "")
	if text == "" {
		return fmt.Errorf("invalid text")
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err := t.dataAPI.Add(
		schemas.DataResponse{
			Data: models.TextData{
				Text: text,
			},
			Meta: meta,
			Type: data.TextDataType,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new text")
	return nil
}

// update обновляет существующую текстовую запись.
//
// args: Аргументы подкоманды, где args[0] - ID записи.
//
// Запрашивает у пользователя новый текст и метаданные, и обновляет запись через API.
// Возвращает ошибку, если не указан ID записи, не удалось получить данные или обновить запись.
func (t *TextCMD) update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to update text must be present")
	}
	id := args[0]

	text := console.GetInput("Enter new text: ", "")
	if text == "" {
		return fmt.Errorf("invalid text")
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err := t.dataAPI.Update(
		schemas.DataResponse{
			Data: models.TextData{
				Text: text,
			},
			Meta: meta,
			Type: data.TextDataType,
			ID:   id,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully updated text")
	return nil
}

// get получает текстовую запись по ID.
//
// args: Аргументы подкоманды, где args[0] - ID записи.
//
// Выводит информацию о записи (ID, текст, метаданные) в консоль.
// Возвращает ошибку, если не указан ID записи или не удалось получить данные.
func (t *TextCMD) get(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to get text must be present")
	}
	dataResponse, err := t.dataAPI.GetByID(args[0], data.TextDataType)
	if err != nil {
		return err
	}
	textData := dataResponse.Data.(models.TextData)
	fmt.Println(
		fmt.Sprintf(
			"ID: %s | Text: %s | Meta: %s",
			dataResponse.ID,
			textData.Text,
			dataResponse.Meta,
		),
	)
	return nil
}

// getAll получает все текстовые записи.
//
// Выводит информацию о всех записях (ID, текст, метаданные) в консоль.
// Возвращает ошибку, если не удалось получить данные.
func (t *TextCMD) getAll(_ []string) error {
	dataResponses, err := t.dataAPI.GetAll(data.TextDataType)
	if err != nil {
		return err
	}
	for _, textData := range dataResponses {
		text := textData.Data.(models.TextData)
		fmt.Println(
			fmt.Sprintf(
				"ID: %s | Text: %s | Meta: %s",
				textData.ID,
				text.Text,
				textData.Meta,
			),
		)
	}
	return nil
}

// delete удаляет текстовую запись по ID.
//
// args: Аргументы подкоманды, где args[0] - ID записи.
//
// Удаляет запись через API.
// Возвращает ошибку, если не указан ID записи или не удалось удалить запись.
func (t *TextCMD) delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to delete text must be present")
	}

	err := t.dataAPI.DeleteByID(args[0], data.TextDataType)
	if err != nil {
		return err
	}
	fmt.Println("Successfully deleted text")
	return nil
}
