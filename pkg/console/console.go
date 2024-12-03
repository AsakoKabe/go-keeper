package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var consoleReader *bufio.Reader

func init() {
	consoleReader = bufio.NewReader(os.Stdin)
}

// GetInput считывает строку ввода из консоли с заданным placeholder'ом и значением по умолчанию.
//
// placeholder: Текст-подсказка, отображаемый в консоли.
// defaultValue: Значение, которое будет возвращено, если пользователь введет пустую строку.
//
// Возвращает введенную пользователем строку или defaultValue, если ввод пустой или произошла ошибка чтения.
func GetInput(placeholder string, defaultValue string) string {
	fmt.Print(placeholder)

	text, err := consoleReader.ReadString('\n')
	if err != nil {
		return defaultValue
	}

	text = strings.Trim(text, "\n\r")

	if text == "" {
		return defaultValue
	}

	return text
}
