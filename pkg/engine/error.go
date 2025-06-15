package engine

import "fmt"

type EngineCompileError struct {
	Message string
}

func (e *EngineCompileError) Error() string {
	return fmt.Sprintf("[engine_compile] error: %s", e.Message)
}

type EngineManagerError struct {
	Message string
}

func (e *EngineManagerError) Error() string {
	return fmt.Sprintf("[engine_manager] error: %s", e.Message)
}

type EngineExecuteError struct {
	Message string
}

func (e *EngineExecuteError) Error() string {
	return fmt.Sprintf("[engine_execute] error: %s", e.Message)
}
