package mailer

import (
	"context"
	"log"
)

type Mailer interface {
	SendVerificationCode(ctx context.Context, email, scene, code string) error
}

type LogMailer struct{}

func NewLogMailer() Mailer {
	return &LogMailer{}
}

func (m *LogMailer) SendVerificationCode(_ context.Context, email, scene, code string) error {
	log.Printf("[mailer] send verification code: email=%s scene=%s code=%s", email, scene, code)
	return nil
}
