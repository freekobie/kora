package main

import (
	"os"

	"github.com/freekobie/kora/mail"
)

type Config struct {
	MailConfig    *mail.Config
	PostgresURL   string
	ServerAddress string
	GCSBucket     string
}

func loadConfig() *Config {

	mailCfg := &mail.Config{
		Host:        os.Getenv("MAIL_HOST"),
		Token:       os.Getenv("MAIL_TOKEN"),
		SenderEmail: os.Getenv("SENDER_EMAIL"),
		SenderName:  os.Getenv("SENDER_NAME"),
	}

	return &Config{
		MailConfig:    mailCfg,
		PostgresURL:   os.Getenv("DB_URL"),
		ServerAddress: os.Getenv("PORT"),
		GCSBucket:     os.Getenv("GCS_BUCKET"),
	}
}
