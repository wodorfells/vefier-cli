package tools

import (
	"fmt"
	"os"
)

type PermMode string

const (
	ModeAskAll       PermMode = "ask_all"
	ModeAllowSession PermMode = "allow_session"
	ModeReadOnly     PermMode = "read_only"
)

type Tool struct {
	Name        string
	Description string
	IsSafe      bool 
	Execute     func(args map[string]string) (string, error)
}

type Registry struct {
	Mode  PermMode
	Tools map[string]Tool
}

func NewRegistry() *Registry {
	r := &Registry{
		Mode:  ModeAskAll,
		Tools: make(map[string]Tool),
	}
	
	r.Tools["read_file"] = Tool{
		Name:        "read_file",
		Description: "Чтение содержимого файла",
		IsSafe:      true, // Read operations are generally safe
		Execute: func(args map[string]string) (string, error) {
			path := args["path"]
			if path == "" {
				return "", fmt.Errorf("path arg required")
			}
			bytes, err := os.ReadFile(path)
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		},
	}
	
	r.Tools["bash"] = Tool{
		Name:        "bash",
		Description: "Выполнение shell-команды",
		IsSafe:      false,
		Execute: func(args map[string]string) (string, error) {
			return "execution simulated (permissions missing)", nil
		},
	}
	
	return r
}

func (r *Registry) ToggleMode() {
	switch r.Mode {
	case ModeAskAll:
		r.Mode = ModeReadOnly
	case ModeReadOnly:
		r.Mode = ModeAllowSession
	case ModeAllowSession:
		r.Mode = ModeAskAll
	}
}

