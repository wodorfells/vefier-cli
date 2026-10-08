package providers

type Provider struct {
	ID       string
	Name     string
	BaseURL  string
	AuthType string // Bearer, X-API-Key
}

func GetCatalog() []Provider {
	return []Provider{
		{ID: "openai", Name: "OpenAI", BaseURL: "https://api.openai.com/v1", AuthType: "Bearer"},
		{ID: "anthropic", Name: "Anthropic", BaseURL: "https://api.anthropic.com/v1", AuthType: "x-api-key"},
		{ID: "gemini", Name: "Google Gemini", BaseURL: "https://generativelanguage.googleapis.com", AuthType: "x-goog-api-key"},
		{ID: "ollama", Name: "Ollama (Local)", BaseURL: "http://localhost:11434/v1", AuthType: "Bearer"},
		{ID: "openrouter", Name: "OpenRouter", BaseURL: "https://openrouter.ai/api/v1", AuthType: "Bearer"},
        {ID: "tako", Name: "Tako", BaseURL: "https://tako.com/api/v3", AuthType: "X-API-Key"},
	}
}