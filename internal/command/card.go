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
	CardCreate = "create"
	CardGet    = "get"
	CardGetAll = "get-all"
	CardUpdate = "update"
	CardDelete = "delete"
)

// CardCMD Команды для управления данными карт.
type CardCMD struct {
	subCommandStrategies map[string]subCommandStrategy
	dataAPI              *api.DataAPI
}

// NewCardCMD создает новый экземпляр CardCMD.
//
// dataAPI: API для взаимодействия с данными.
func NewCardCMD(dataAPI *api.DataAPI) *CardCMD {
	cardCMD := &CardCMD{dataAPI: dataAPI}
	cardCMD.subCommandStrategies = map[string]subCommandStrategy{
		CardCreate: cardCMD.create,
		CardGet:    cardCMD.get,
		CardUpdate: cardCMD.update,
		CardDelete: cardCMD.delete,
		CardGetAll: cardCMD.getAll,
	}
	return cardCMD
}

// Execute выполняет подкоманду для управления картами.
//
// args: Аргументы команды, где args[0] - подкоманда (create, get, update, delete, get-all).
//
// Возвращает ошибку, если подкоманда не найдена или произошла ошибка при ее выполнении.
func (c *CardCMD) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("sub command for card not found")
	}

	subCMD, ok := c.subCommandStrategies[args[0]]
	if !ok {
		return fmt.Errorf("sub command not found")
	}

	return subCMD(args[1:])
}

// create создает новую карту.
//
// _: Аргументы подкоманды (не используются).
//
// Запрашивает у пользователя данные карты (номер, срок действия, CVV, метаданные) и сохраняет их через API.
// Возвращает ошибку, если не удалось получить данные или сохранить карту.
func (c *CardCMD) create(_ []string) error {
	number := console.GetInput("Enter card number: ", "")
	if number == "" {
		return fmt.Errorf("invalid card number")
	}

	expiredAt := console.GetInput("Enter expiration date (MM/YY): ", "")
	if expiredAt == "" {
		return fmt.Errorf("invalid expiration date")
	}

	cvv := console.GetInput("Enter CVV: ", "")
	if cvv == "" {
		return fmt.Errorf("invalid CVV")
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err := c.dataAPI.Add(
		schemas.DataResponse{
			Data: models.CardData{
				Number:    number,
				ExpiredAt: expiredAt,
				CVV:       cvv,
			},
			Meta: meta,
			Type: data.CardDataType,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully saved new card")
	return nil
}

// update обновляет существующую карту.
//
// args: Аргументы подкоманды, где args[0] - ID карты.
//
// Запрашивает у пользователя данные карты (номер, срок действия, CVV, метаданные) и обновляет их через API.
// Возвращает ошибку, если не указан ID карты, не удалось получить данные или обновить карту.
func (c *CardCMD) update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to update card must be present")
	}
	id := args[0]

	number := console.GetInput("Enter card number: ", "")
	if number == "" {
		return fmt.Errorf("invalid card number")
	}

	expiredAt := console.GetInput("Enter expiration date (MM/YY): ", "")
	if expiredAt == "" {
		return fmt.Errorf("invalid expiration date")
	}

	cvv := console.GetInput("Enter CVV: ", "")
	if cvv == "" {
		return fmt.Errorf("invalid CVV")
	}

	meta := console.GetInput("Enter meta: ", "")
	_, err := c.dataAPI.Update(
		schemas.DataResponse{
			Data: models.CardData{
				Number:    number,
				ExpiredAt: expiredAt,
				CVV:       cvv,
			},
			Meta: meta,
			Type: data.CardDataType,
			ID:   id,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Successfully updated card")
	return nil
}

// get получает информацию о карте по ID.
//
// args: Аргументы подкоманды, где args[0] - ID карты.
//
// Выводит информацию о карте в консоль.
// Возвращает ошибку, если не указан ID карты или не удалось получить данные о карте.
func (c *CardCMD) get(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to get card must be present")
	}
	dataResponse, err := c.dataAPI.GetByID(args[0], data.CardDataType)
	if err != nil {
		return err
	}
	cardData := dataResponse.Data.(models.CardData)
	fmt.Println(
		fmt.Sprintf(
			"ID: %s | Number: %s | Expiration Date: %s | CVV: %s | Meta: %s",
			dataResponse.ID,
			cardData.Number,
			cardData.ExpiredAt,
			cardData.CVV,
			dataResponse.Meta,
		),
	)
	return nil
}

// getAll получает информацию о всех картах.
//
// _: Аргументы подкоманды (не используются).
//
// Выводит информацию о всех картах в консоль.
// Возвращает ошибку, если не удалось получить данные о картах.
func (c *CardCMD) getAll(_ []string) error {
	dataResponses, err := c.dataAPI.GetAll(data.CardDataType)
	if err != nil {
		return err
	}
	for _, cardData := range dataResponses {
		card := cardData.Data.(models.CardData)
		fmt.Println(
			fmt.Sprintf(
				"ID: %s | Number: %s | Expiration Date: %s | CVV: %s | Meta: %s",
				cardData.ID,
				card.Number,
				card.ExpiredAt,
				card.CVV,
				cardData.Meta,
			),
		)
	}
	return nil
}

// delete удаляет карту по ID.
//
// args: Аргументы подкоманды, где args[0] - ID карты.
//
// Удаляет карту через API.
// Возвращает ошибку, если не указан ID карты или не удалось удалить карту.
func (c *CardCMD) delete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("id to delete card must be present")
	}

	err := c.dataAPI.DeleteByID(args[0], data.CardDataType)
	if err != nil {
		return err
	}
	fmt.Println("Successfully deleted card")
	return nil
}
