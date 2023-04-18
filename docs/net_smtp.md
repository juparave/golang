# Senging mail using `net/smtp` package

The `net/smtp` package in Go provides a simple way to send email messages over
the Simple Mail Transfer Protocol (SMTP). This package is included in the
standard library, so no external dependencies are required.

## Basic usage

Here is a basic example of how to send an email using the `net/smtp` package:

```go
package main

import (
    "net/smtp"
)

func main() {
    // Set up SMTP connection
    smtpHost := "smtp.example.com"
    smtpPort := "587"
    auth := smtp.PlainAuth("", "user@example.com", "password", smtpHost)

    // Set up email headers and content
    from := "user@example.com"
    to := []string{"recipient@example.com"}
    subject := "Subject: Hello from Go\n"
    body := "This is a test email sent from Go."

    // Send the email
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(subject+body))
    if err != nil {
        panic(err)
    }
}
```

In this example, we first set up the SMTP connection by specifying the SMTP
server hostname, port number, and authentication credentials. We then set up
the email headers and content by specifying the sender, recipient, subject, and
body. Finally, we use the `smtp.SendMail` function to send the email.

## Sending to Multiple Recipients

To send an email to multiple recipients, simply pass a slice of email addresses
to the to parameter of the smtp.SendMail function:

```go
to := []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"}
err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(subject+body))
```

Avoid FRC-5321 Errors

```go
to := []string{" recipient1@example.com ", "recipient2@example.com", " recipient3@example.com "}
for i, addr := range to {
    to[i] = strings.TrimSpace(addr)
}
err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(subject+body))
```

## Formatting the Email Body

The `smtp.SendMail` function expects the email body to be in the format of a byte
slice containing the email headers followed by a blank line and the email
content. To send an HTML email, include the appropriate `Content-Type` header in
the email headers:

```go
subject := "Subject: Hello from Go\n"
body := "<html><body><h1>This is an HTML email sent from Go.</h1></body></html>"
headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nContent-Type: text/html\r\nSubject: %s\r\n\r\n", from, strings.Join(to, ","), subject)
err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(headers+body))
```

If rendering an html template:

```go
// hold the body of the email
var body bytes.Buffer

// Parse the html template
tmpl, err := template.ParseFS(templates, "templates/"+m.TemplateName+".html")
if err != nil {
    app.ErrorLog.Println("Error parsing template:", err)
    return err
}

// Render the template with the given data
err = tmpl.Execute(&body, m.TemplateCtx)
if err != nil {
    app.ErrorLog.Println("Error executing template:", err)
    return err
}

from := "user@example.com"
fromHeader := fmt.Sprintf("From: Example User <%s>\n", from)
to := []string{"recipient@example.com"}
for i, addr := range to {
    // clear whitespaces from to's list
    to[i] = strings.TrimSpace(addr)
}
toHeader := fmt.Sprintf("To: %s\n", strings.Join(to, ","))
subject := "Hello from go"
subjectHeader := "Subject: " + subject + "\n"
mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
msg := []byte(fromHeader + toHeader + subject + mime + body.String())

// Send the email
err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
if err != nil {
    app.ErrorLog.Println("Error sending email:", err)
    return err
}
```


## Email Attachments

In order to send an email with attachments, you need to create a MIME multipart
message. Here's an example code that shows how to send an email with a PDF file
as an attachment:

```go
// Set up SMTP connection
smtpHost := app.Email.Host
smtpPort := app.Email.Port
auth := smtp.PlainAuth("",
    app.Email.Account,
    app.Email.Password,
    smtpHost)

// Create a new buffer for the email
var emailBuffer bytes.Buffer

// Create a multipart writer
multipartWriter := multipart.NewWriter(&emailBuffer)

// Set up email headers and content
from := app.Email.Account
fromHeader := fmt.Sprintf("From: Example User <%s>\r\n", from)
// Set the toAddresses address
toAddresses := strings.Split(m.To, ",")
for i, address := range toAddresses {
    toAddresses[i] = strings.TrimSpace(address)
}
toHeader := fmt.Sprintf("To: %s\r\n", m.To)
// Email Subject
subject := "Subject: " + m.Subject + "\r\n"
// Set mime header
mimeHeader := "MIME-version: 1.0;\r\nContent-Type: " + multipartWriter.FormDataContentType() + "\r\n\r\n"
// Concatenate to create mail headers
headers := fromHeader + toHeader + subject + mimeHeader

// write the headers to the buffer
emailBuffer.WriteString(headers)

// Create a new part for the html body
htmlPart, err := multipartWriter.CreatePart(
    textproto.MIMEHeader{
        "Content-Type":        []string{"text/html; charset=UTF-8"},
        "Content-Disposition": {"inline"},
    })
if err != nil {
    app.ErrorLog.Println("Error creating html part:", err)
    return err
}

// Register the custom template function before parsing the template
tmpl := template.New(fmt.Sprintf("%s.html", m.TemplateName)).Funcs(template.FuncMap{
    "safeHTML": safeHTML,
})

// Render HTML template
tmpl, err = tmpl.ParseFS(templates, "templates/"+m.TemplateName+".html")
if err != nil {
    app.ErrorLog.Println("Error parsing template:", err)
    return err
}

// Render the template with the given data
err = tmpl.Execute(htmlPart, m.TemplateCtx)
if err != nil {
    app.ErrorLog.Println("Error executing template:", err)
    return err
}

// Check if there's an attachment and add it to the email
if *m.Attachment != "" {
    // Decode the base64 encoded attachment
    attachmentData, err := base64.StdEncoding.DecodeString(*m.Attachment)
    if err != nil {
        app.ErrorLog.Println("Error decoding attachment:", err)
        return err
    }

    // set the filename that appears in the mail
    filename := "your_filename.pdf"
    attachmentPart, err := multipartWriter.CreatePart(textproto.MIMEHeader{
        "Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filename)},
        "Content-Type":              {"application/pdf"}, // Set the appropriate content type based on the file
        "Content-Transfer-Encoding": {"base64"},
    })
    if err != nil {
        app.ErrorLog.Println("Error creating attachment part:", err)
        return err
    }

    // Copy the attachment to the email
    attachmentBuffer := bytes.NewReader(attachmentData)
    _, err = io.Copy(attachmentPart, attachmentBuffer)
    if err != nil {
        app.ErrorLog.Println("Error copying attachment to email:", err)
        return err
    }
    // add the attachment to the email
    _, err = emailBuffer.Write(attachmentData)
    if err != nil {
        app.ErrorLog.Println("Error writing attachment to email:", err)
        return err
    }
}

// Close the multipart writer
err = multipartWriter.Close()
if err != nil {
    app.ErrorLog.Println("Error closing multipart writer:", err)
    return err
}

// Send the email
err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toAddresses, emailBuffer.Bytes())
if err != nil {
    app.ErrorLog.Println("Error sending email:", err)
    return err
}

```
