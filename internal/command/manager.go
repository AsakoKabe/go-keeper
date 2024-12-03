package command

// Manager управляет доступными командами.
type Manager struct {
	commands map[string]Command
}

// NewManager создает новый экземпляр Manager.
func NewManager() *Manager {
	return &Manager{
		commands: make(map[string]Command),
	}
}

// RunCommand запускает команду с указанным именем и аргументами.
//
// cmd: Имя команды для запуска.
// args: Аргументы команды.
//
// Возвращает ErrCommandNotFound, если команда не найдена, или ошибку, возвращенную методом Execute команды.
func (m *Manager) RunCommand(cmd string, args []string) error {
	command, ok := m.commands[cmd]
	if !ok {
		return ErrCommandNotFound
	}

	return command.Execute(args)
}

// AddCommand добавляет команду в менеджер.
//
// name: Имя команды.
// command: Экземпляр команды.
func (m *Manager) AddCommand(name string, command Command) {
	m.commands[name] = command
}
