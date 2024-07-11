package mailer

import (
	"fmt"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"go.uber.org/zap"
	"net/smtp"
	"strings"
)

type Mailer struct {
	c *config.Config
	l *logger.Logger
}

func (m *Mailer) Send(to, topic, message string) {
	body := []byte(strings.Join([]string{
		"Subject: " + topic,
		"User: " + "Nightell" + " <" + m.c.Mailer.Username + ">",
		"To: " + to,
		"\r\n" + message,
	}, "\r\n"))

	auth := smtp.PlainAuth("", m.c.Mailer.Username, m.c.Mailer.Password, m.c.Mailer.SmtpServer)
	server := m.c.Mailer.SmtpServer + ":" + fmt.Sprintf("%d", m.c.Mailer.SmtpPort)

	err := smtp.SendMail(server, auth, m.c.Mailer.Username, []string{to}, body)
	if err != nil {
		m.l.Info("mailer: failed", zap.String("to", to), zap.Error(err))
		return
	}

	m.l.Info("mailer: sent successfully", zap.String("to", to))
}

func (m *Mailer) SendWelcome(to, username string) {
	message := strings.Join([]string{
		"Dear " + username + ",",
		"Congratulations on successfully registering with Nightell!",
		"You can now sign in using our app and share your voice with the world...",
		"https://nightell.neatplex.com",
	}, "\r\n")
	m.Send(to, "Welcome to Nightell!", message)
}

func (m *Mailer) SendDeleteAccount(to, username, link string) {
	message := strings.Join([]string{
		"Dear " + username + ",",
		"We have received a request to delete your account associated with this email address.",
		"Please note that deleting your account is a permanent action and cannot be undone. " +
			"All your data, including any saved preferences and history, will be permanently erased.",
		"If you did not make this request, please ignore this email. Your account will remain unchanged.",
		"To confirm the deletion of your account, please click the link below:",
		link,
		"",
		"If you have any questions or need assistance, please contact our support team at nightell@neatplex.com.",
		"Thank you for being a part of our community.",
		"",
		"https://nightell.neatplex.com.",
	}, "\r\n")
	m.Send(to, "Nightell: Delete Account", message)
}

func New(c *config.Config, l *logger.Logger) *Mailer {
	return &Mailer{
		l: l,
		c: c,
	}
}
