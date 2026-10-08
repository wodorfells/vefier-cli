package keys

import (
	"bufio"
	"strings"
)

// ParseKeys reads plain text with multiple formats
func ParseKeys(data string) []Key {
	var imported []Key
	scanner := bufio.NewScanner(strings.NewReader(data))
	
	count := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		name := ""
		value := line
		
		if idx := strings.IndexAny(line, "=:"); idx != -1 {
			name = strings.TrimSpace(line[:idx])
			value = strings.TrimSpace(line[idx+1:])
		}
		
		if name == "" {
			name = "Key " + string(rune(count+'0'))
			count++
		}
		
		provider := detectProvider(value)
		imported = append(imported, Key{
			ID:       name, 
			Name:     name,
			Value:    value,
			Provider: provider,
			Status:   "untested",
		})
	}
	return imported
}

func detectProvider(val string) string {
	if strings.HasPrefix(val, "sk-ant") {
		return "anthropic"
	} else if strings.HasPrefix(val, "tako_sk") {
		return "tako" 
	} else if strings.HasPrefix(val, "sk-") {
		return "openai" 
	}
	return "openai"
}

