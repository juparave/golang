## HTML templates with go

#### Embed a directory

```go
package transcribe

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*
var templates embed.FS

func createChatDisplay(leftText, rightText string) (string, error) {
	type ChatMessage struct {
		Text   string
		IsLeft bool
	}

	var chatMessages []ChatMessage

	// Parse left text as vtt format and append to chatMessages with IsLeft set to true
	// Parse right text as vtt format and append to chatMessages with IsLeft set to false

	// Render HTML template with chatMessages
	t, err := template.ParseFS(templates, "templates/chat.html")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, chatMessages)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
```
