package mcp

type Client struct {
	Name   string
	Cmd    string
	Status string
}

type Manager struct {
	Clients []Client
}

func NewManager() *Manager {
	return &Manager{
		Clients: []Client{
			{Name: "Filesystem", Cmd: "npx -y @modelcontextprotocol/server-filesystem", Status: "running"},
			{Name: "Postgres", Cmd: "npx -y @modelcontextprotocol/server-postgres", Status: "stopped"},
		},
	}
}