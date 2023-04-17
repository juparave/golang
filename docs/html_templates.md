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

Add functions to parser.

```go
import (
    "html/template"
)

// Define a custom template function to mark content as safe HTML
func safeHTML(s string) template.HTML {
    return template.HTML(s)
}
```

Using embeded FS the function needs to add the function before parsing it.

```go
// Register the custom template function
t := template.New("chat.html").Funcs(template.FuncMap{
    "safeHTML": safeHTML,
})

// Parse your template
t, err := t.ParseFS(templates, "templates/chat.html")
if err != nil {
    // handle error
}
```

Simple html return from template

```go
//go:embed templates/*
var templates embed.FS

// PrivacyPolicy returns the Privacy Policy
func PrivacyPolicy(c *fiber.Ctx) error {
	page, err := templates.ReadFile("templates/privacy-policy.html")
	if err != nil {
		app.ErrorLog.Println(err)
		return err
	}

	// Set appropriate headers for HTML response
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.Send(page)
}
```
