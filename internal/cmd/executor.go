package cmd

var defaultExecutor CommandExecutor = &SimpleCommandExecutor{}

type CommandExecutor interface {
	Execute(cmd Command) (output []byte, err error)
}

func GetCommandExecutor() CommandExecutor {
	return defaultExecutor
}

func SetCommandExecutor(executor CommandExecutor) {
	defaultExecutor = executor
}

type SimpleCommandExecutor struct {
}

func (e *SimpleCommandExecutor) Execute(cmd Command) (output []byte, err error) {
	return cmd.CombinedOutput()
}
