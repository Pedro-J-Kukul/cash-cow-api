// Fileanme: internal/mailer/mailer.go
package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail"
)

//go:embed "templates/*"
var templateFS embed.FS // embed the templates directory

// Mailer struct to hold mailer configuration and state
type Mailer struct {
	dialer *mail.Dialer // SMTP dialer for sending emails
	sender string       // sender email address
}

// New creates a new Mailer instance with the provided SMTP configuration
func New(host string, port int, username, password, sender string) *Mailer {
	dialer := mail.NewDialer(host, port, username, password) // create a new SMTP dialer
	dialer.Timeout = 5 * time.Second

	// return a pointer to a new Mailer instance
	return &Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send sends an email using the Mailer instance
func (m *Mailer) Send(to, templateFile string, data any) error {
	tmpl, err := template.ParseFS(templateFS, "templates/"+templateFile) // parse the email template
	if err != nil {
		return err // return error if template parsing fails
	}

	subject := new(bytes.Buffer)                         // buffer to hold the email subject
	err = tmpl.ExecuteTemplate(subject, "subject", data) // execute the subject template
	if err != nil {
		return err // return error if subject template execution fails
	}

	plainBody := new(bytes.Buffer)                           // buffer to hold the plain text body
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data) // execute the plain body template
	if err != nil {
		return err // return error if plain body template execution fails
	}

	htmlBody := new(bytes.Buffer)                          // buffer to hold the HTML body
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data) // execute the HTML body template
	if err != nil {
		return err // return error if HTML body template execution fails
	}

	msg := mail.NewMessage()                           // create a new email message
	msg.SetHeader("From", m.sender)                    // set the sender header
	msg.SetHeader("To", to)                            // set the recipient header
	msg.SetHeader("Subject", subject.String())         // set the subject header
	msg.SetBody("text/plain", plainBody.String())      // set the plain text body
	msg.AddAlternative("text/html", htmlBody.String()) // add the HTML body as an alternative

	// Sends the email 3 times before returning an error
	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(msg) // attempt to send the email
		if err == nil {
			break // break the loop if sending is successful
		}
		time.Sleep(500 * time.Millisecond) // wait before retrying
	}

	return err // return any error encountered during sending
}
