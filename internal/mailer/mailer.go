package mailer

import (
	"fmt"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"go.uber.org/zap"
	"net/smtp"
)

type Mailer struct {
	c *config.Config
	l *logger.Logger
}

func (m *Mailer) Send(to, message string) {
	body := []byte(
		"Subject: Hello!\r\n" +
			"From: " + m.c.Mailer.Username + "\r\n" +
			"To: " + to + "\r\n" +
			"\r\n" +
			message,
	)

	auth := smtp.PlainAuth("", m.c.Mailer.Username, m.c.Mailer.Password, m.c.Mailer.SmtpServer)
	server := m.c.Mailer.SmtpServer + ":" + fmt.Sprintf("%d", m.c.Mailer.SmtpPort)

	err := smtp.SendMail(server, auth, m.c.Mailer.Username, []string{to}, body)
	if err != nil {
		m.l.Info("mailer: failed", zap.String("to", to), zap.Error(err))
		return
	}

	m.l.Info("mailer: sent successfully", zap.String("to", to))
}

func New(c *config.Config, l *logger.Logger) *Mailer {
	return &Mailer{
		l: l,
		c: c,
	}
}
