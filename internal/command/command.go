package command

// Command интерфейс для выполнения команд.
type Command interface {
	Execute(args []string) error
}

type subCommandStrategy func(args []string) error
