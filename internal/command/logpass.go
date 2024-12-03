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
	LogPassCreate = "create"
	LogPassGet    = "get"
	LogPassGetAll = "get-all"
	LogPassUpdate = "update"
	LogPassDelete = "delete"
)

// LogPassCMD предоставляет команды для управления данными логинов и паролей.
type LogPassCMD struct {
	subCommandStrategies map[string]subCommandStrategy
	dataAPI              *api.DataAPI
}

// NewLogPassCMD создает новый экземпляр LogPassCMD.
//
// dataAPI: API для взаимодействия с данными.
func NewLogPassCMD(dataAPI *api.DataAPI) *LogPassCMD {
	logPass := &LogPassCMD{dataAPI: dataAPI}
	logPass.subCommandStrategies = map[string]subCommandStrategy{
		LogPassCreate: logPass.create,
		LogPassGet:    logPass.get,
		LogPassUpdate: logPass.update,
		LogPassDelete: logPass.delete,
		LogPassGetAll: logPass.getAll,
	}
	return logPass
}

// Execute выполняет подкоманду для управления логинами и паролями.
//
// args: Аргументы команды, где args[0] - подкоманда (create, get, update, delete, get-all).
//
// Возвращает ошибку, если подкоманда не найдена или произошла ошибка при ее выполнении.
func (l *LogPassCMD) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("sub command for log-pass not found")
	}

	subCMD, ok := l.subCommandStrategies[args[0]]
	if !ok {
		return fmt.Errorf("sub command not found")
	}

	return subCMD(args[1:])
}

func (l *LogPassCMD) create(_ []string) error {
	login := console.GetInput("Enter login: ", "")
	if login == "" {
		return fmt.Errorf("invalid login")
	}

	password := console.GetInput("Enter password: ", "")
	if password == "" {
		return fmt.Errorf("invalid password")
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err := l.dataAPI.Add(
		schemas.DataResponse{
			Data: models.LogPassData{
				Login:    login,
				Password: password,
			},
			Meta: meta,
			Type: data.LogPassDataType,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new logpass")
	return nil
}

// update обновляет существующую запись логина и пароля.
//
// args: Аргументы подкоманды, где args[0] - ID записи.
//
// Запрашивает у пользователя логин, пароль и метаданные, и обновляет запись через API.
// Возвращает ошибку, если не указан ID записи, не удалось получить данные или обновить запись.
func (l *LogPassCMD) update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to update logpass must present")
	}
	id := args[0]

	login := console.GetInput("Enter login: ", "")
	if login == "" {
		return fmt.Errorf("invalid login")
	}

	password := console.GetInput("Enter password: ", "")
	if password == "" {
		return fmt.Errorf("invalid password")
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err := l.dataAPI.Update(
		schemas.DataResponse{
			Data: models.LogPassData{
				Login:    login,
				Password: password,
			},
			Meta: meta,
			Type: data.LogPassDataType,
			ID:   id,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully updates logpass")
	return nil
}

// get получает запись логина и пароля по ID.
//
// args: Аргументы подкоманды, где args[0] - ID записи.
//
// Выводит информацию о записи (ID, логин, пароль, метаданные) в консоль.
// Возвращает ошибку, если не указан ID записи или не удалось получить данные.
func (l *LogPassCMD) get(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to get logpass must present")
	}
	dataResponse, err := l.dataAPI.GetByID(args[0], data.LogPassDataType)
	if err != nil {
		return err
	}
	logPassData := dataResponse.Data.(models.LogPassData)
	fmt.Println(
		fmt.Sprintf(
			"ID: %s | Login: %s Password: %s | Meta: %s",
			dataResponse.ID,
			logPassData.Login,
			logPassData.Password,
			dataResponse.Meta,
		),
	)
	return nil
}

// getAll получает все записи логинов и паролей.
//
// _: Аргументы подкоманды (не используются).
//
// Выводит информацию о всех записях (ID, логин, пароль, метаданные) в консоль.
// Возвращает ошибку, если не удалось получить данные.
func (l *LogPassCMD) getAll(_ []string) error {
	dataResponses, err := l.dataAPI.GetAll(data.LogPassDataType)
	if err != nil {
		return err
	}
	for _, userData := range dataResponses {
		logPassData := userData.Data.(models.LogPassData)
		fmt.Println(
			fmt.Sprintf(
				"ID: %s | Login: %s Password: %s | Meta: %s",
				userData.ID,
				logPassData.Login,
				logPassData.Password,
				userData.Meta,
			),
		)

	}
	return nil
}

// delete удаляет запись логина и пароля по ID.
//
// args: Аргументы подкоманды, где args[0] - ID записи.
//
// Удаляет запись через API.
// Возвращает ошибку, если не указан ID записи или не удалось удалить запись.
func (l *LogPassCMD) delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to delete logpass must present")
	}

	err := l.dataAPI.DeleteByID(args[0], data.LogPassDataType)
	if err != nil {
		return err
	}
	fmt.Println("successfully deleted")
	return nil
}
