package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wodorfells/vefier-cli/internal/providers"
)

type StreamChunk struct {
	Content string
	Error   error
	Done    bool
}

func StreamCompletion(ctx context.Context, p providers.Provider, token string, model string, prompt string, ch chan<- StreamChunk) {
	defer close(ch)
	
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
		"stream":   true,
	})

	url := p.BaseURL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "chat/completions" // Base assumption for OpenAI-compatible

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		ch <- StreamChunk{Error: err}
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	if strings.ToLower(p.AuthType) == "bearer" || strings.ToLower(p.AuthType) == "openai" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.Header.Set(p.AuthType, token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- StreamChunk{Error: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ch <- StreamChunk{Error: fmt.Errorf("API error: %s", resp.Status)}
		return // Here you would normally trigger key rotation (401, 429)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, []byte("data: ")) {
			continue
		}
		
		data := bytes.TrimPrefix(line, []byte("data: "))
		if string(data) == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(data, &chunk); err == nil {
			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				if content != "" {
					ch <- StreamChunk{Content: content}
				}
			}
		}
	}
	ch <- StreamChunk{Done: true}
}

