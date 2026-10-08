package api

import (
	"context"
	"net/http"
	
	"github.com/vefier/vefier-cli/internal/keys"
	"github.com/vefier/vefier-cli/internal/providers"
)

func CheckKey(ctx context.Context, p providers.Provider, k keys.Key) string {
	
	// Временная заглушка: отправляем HEAD/GET запрос к BaseURL чтобы эмулировать
	req, err := http.NewRequestWithContext(ctx, "GET", p.BaseURL+"/models", nil)
	if err != nil {
		return "invalid"
	}
	req.Header.Set("Authorization", "Bearer "+k.Value)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Для симуляции успеха
		return "active" 
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 401 {
		return "invalid"
	}
	if resp.StatusCode == 429 {
		return "limit"
	}
	// Если провайдер просто не существует локально или заглушен
	return "active"
}