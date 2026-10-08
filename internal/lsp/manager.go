package lsp

type Client struct {
	Language string
	Cmd      string
	Status   string
}

type Manager struct {
	Clients []Client
}

func NewManager() *Manager {
	return &Manager{
		Clients: []Client{
			{Language: "go", Cmd: "gopls", Status: "running"},
			{Language: "python", Cmd: "pyright", Status: "stopped"},
		},
	}
}