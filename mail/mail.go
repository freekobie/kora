package mail

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"
)

//go:embed templates
var mailFS embed.FS

type Mailer struct {
	cfg *Config
}

type Config struct {
	Host        string
	Token       string
	Timeout     time.Duration
	SenderName  string
	SenderEmail string
}

type Address struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Message struct {
	From    Address   `json:"from"`
	To      []Address `json:"to"`
	Subject string    `json:"subject"`
	Text    string    `json:"text"`
	HTML    string    `json:"html"`
}

type Data struct {
	Address Address
	Code    any
}

func NewMailer(config *Config) *Mailer {
	mailer := &Mailer{
		cfg: config,
	}
	return mailer
}

func (m *Mailer) generateMessage(recipient []Address, templateFile string, data any) (*Message, error) {
	tmpl, err := template.ParseFS(mailFS, fmt.Sprintf("template/%s", templateFile))
	if err != nil {
		slog.Error("error parsing FS", "error", err)
		return nil, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		slog.Error("error executing subject template", "error", err)
		return nil, err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "text", data)
	if err != nil {
		slog.Error("error executing text template", "error", err)
		return nil, err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "html", data)
	if err != nil {
		slog.Error("error executing html template", "error", err)
		return nil, err
	}

	msg := Message{
		From:    Address{Email: m.cfg.SenderEmail, Name: m.cfg.SenderName},
		To:      recipient,
		Subject: subject.String(),
		Text:    plainBody.String(),
		HTML:    htmlBody.String(),
	}

	return &msg, nil
}

func (m *Mailer) Send(recipients []Address, templateFile string, data any) error {
	msg, err := m.generateMessage(recipients, templateFile, data)
	if err != nil {
		return err
	}

	msgJson, err := json.Marshal(msg)
	if err != nil {
		slog.Error("error marshalling mail message", "error", err)
		return err
	}

	req, err := http.NewRequest("POST", m.cfg.Host, bytes.NewBuffer(msgJson))
	if err != nil {
		slog.Error("error creating request", "error", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.cfg.Token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: m.cfg.Timeout}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("error sending request", "error", err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("error reading body", "error", err)
			return err
		}

		slog.Error("error executing template", slog.String("status", res.Status), slog.Any("body", body))
		return errors.New("error sending email: " + res.Status)
	}

	return nil
}
