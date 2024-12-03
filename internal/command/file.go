package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"go-keeper/internal/api"
	"go-keeper/internal/http/rest/schemas"
	"go-keeper/internal/keeper/data"
	"go-keeper/internal/models"
	"go-keeper/pkg/console"
)

const (
	FileCreate = "create"
	FileGet    = "get"
	FileGetAll = "get-all"
	FileUpdate = "update"
	FileDelete = "delete"
)

// FileCMD команд для управления файлом
type FileCMD struct {
	subCommandStrategies map[string]subCommandStrategy
	dataAPI              *api.DataAPI
}

// NewFileCMD создает новый экземпляр FileCMD.
//
// dataAPI: API для взаимодействия с данными.
func NewFileCMD(dataAPI *api.DataAPI) *FileCMD {
	fileCMD := &FileCMD{dataAPI: dataAPI}
	fileCMD.subCommandStrategies = map[string]subCommandStrategy{
		FileCreate: fileCMD.create,
		FileGet:    fileCMD.get,
		FileUpdate: fileCMD.update,
		FileDelete: fileCMD.delete,
		FileGetAll: fileCMD.getAll,
	}
	return fileCMD
}

// Execute выполняет подкоманду для управления файлами.
//
// args: Аргументы команды, где args[0] - подкоманда (create, get, update, delete, get-all).
//
// Возвращает ошибку, если подкоманда не найдена или произошла ошибка при ее выполнении.
func (f *FileCMD) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("sub command for file not found")
	}

	subCMD, ok := f.subCommandStrategies[args[0]]
	if !ok {
		return fmt.Errorf("sub command not found")
	}

	return subCMD(args[1:])
}

// create создает новый файл.
//
// _: Аргументы подкоманды (не используются).
//
// Запрашивает у пользователя путь к файлу и метаданные, считывает содержимое файла и сохраняет его через API.
// Возвращает ошибку, если не удалось получить путь к файлу, прочитать файл или сохранить данные.
func (f *FileCMD) create(_ []string) error {
	filePath := console.GetInput("Enter file path: ", "")
	if filePath == "" {
		return fmt.Errorf("invalid file path")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err = f.dataAPI.Add(
		schemas.DataResponse{
			Data: models.FileData{
				Content: content,
				Name:    filePath,
			},
			Meta: meta,
			Type: data.FileDataType,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new file")
	return nil
}

// update обновляет существующий файл.
//
// args: Аргументы подкоманды, где args[0] - ID файла.
//
// Запрашивает у пользователя новый путь к файлу и метаданные, считывает содержимое файла и обновляет данные через API.
// Возвращает ошибку, если не указан ID файла, не удалось получить путь к файлу, прочитать файл или обновить данные.
func (f *FileCMD) update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to update file must be present")
	}
	id := args[0]

	filePath := console.GetInput("Enter new file path: ", "")
	if filePath == "" {
		return fmt.Errorf("invalid file path")
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err = f.dataAPI.Update(
		schemas.DataResponse{
			Data: models.FileData{
				Content: content,
				Name:    filePath,
			},
			Meta: meta,
			Type: data.FileDataType,
			ID:   id,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully updated file")
	return nil
}

// get получает файл по ID.
//
// args: Аргументы подкоманды, где args[0] - ID файла.
//
// Получает данные файла через API и сохраняет их в файл с указанным именем.
// Возвращает ошибку, если не указан ID файла, не удалось получить данные
func (f *FileCMD) get(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to get file must be present")
	}
	dataResponse, err := f.dataAPI.GetByID(args[0], data.FileDataType)
	if err != nil {
		return err
	}
	fileData := dataResponse.Data.(models.FileData)

	fmt.Println(
		fmt.Sprintf(
			"ID: %s | Name: %s | Meta: %s",
			dataResponse.ID,
			fileData.Name,
			dataResponse.Meta,
		),
	)
	return nil
}

// getAll получает все файлы.
//
// _: Аргументы подкоманды (не используются).
//
// Получает данные всех файлов через API и сохраняет каждый файл с соответствующим именем.
// Выводит информацию о каждом файле (ID, имя, метаданные).
func (f *FileCMD) getAll(_ []string) error {
	dataResponses, err := f.dataAPI.GetAll(data.FileDataType)
	if err != nil {
		return err
	}
	for _, fileData := range dataResponses {
		file := fileData.Data.(models.FileData)

		fmt.Println(
			fmt.Sprintf(
				"ID: %s | Name: %s | Meta: %s",
				fileData.ID,
				file.Name,
				fileData.Meta,
			),
		)
	}
	return nil
}

// delete удаляет файл по ID.
//
// args: Аргументы подкоманды, где args[0] - ID файла.
//
// Удаляет файл через API.
// Возвращает ошибку, если не указан ID файла
func (f *FileCMD) delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to delete file must be present")
	}

	err := f.dataAPI.DeleteByID(args[0], data.FileDataType)
	if err != nil {
		return err
	}
	fmt.Println("Successfully deleted file")
	return nil
}
