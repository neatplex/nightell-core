package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"go.uber.org/zap"
)

const (
	defaultMailerEndpoint = "https://api.zeptomail.eu/v1.1/email"

	driverSMTP      = "smtp"
	driverZeptoMail = "zeptomail"

	encryptionTLS = "tls"
)

type Mailer struct {
	c       *config.Config
	l       *logger.Logger
	client  *http.Client
	timeout time.Duration
}

type zeptoAddress struct {
	Address string `json:"address"`
}

type zeptoRecipient struct {
	EmailAddress zeptoAddress `json:"email_address"`
}

type zeptoMailRequest struct {
	From     zeptoAddress     `json:"from"`
	To       []zeptoRecipient `json:"to"`
	Subject  string           `json:"subject"`
	HTMLBody string           `json:"htmlbody"`
}

func (m *Mailer) Send(to, topic, message string) {
	if strings.TrimSpace(m.c.Mailer.Sender) == "" {
		m.l.Error("mailer: sender email is not configured")
		return
	}

	var err error
	switch m.c.Mailer.Driver {
	case driverSMTP:
		err = m.sendSMTP(to, topic, message)
	case driverZeptoMail, "":
		err = m.sendZeptoMail(to, topic, message)
	default:
		err = fmt.Errorf("unknown driver %q", m.c.Mailer.Driver)
	}
	if err != nil {
		m.l.Error("mailer: failed", zap.String("to", to), zap.Error(err))
		return
	}

	m.l.Info("mailer: sent successfully", zap.String("to", to))
}

func (m *Mailer) sendZeptoMail(to, topic, message string) error {
	if m.client == nil {
		return errors.New("http client not initialized")
	}

	if strings.TrimSpace(m.c.Mailer.APIKey) == "" {
		return errors.New("api key is not configured")
	}

	endpoint := m.c.Mailer.APIEndpoint
	if strings.TrimSpace(endpoint) == "" {
		endpoint = defaultMailerEndpoint
	}

	requestBody := zeptoMailRequest{
		From: zeptoAddress{
			Address: m.c.Mailer.Sender,
		},
		To: []zeptoRecipient{
			{EmailAddress: zeptoAddress{Address: to}},
		},
		Subject:  topic,
		HTMLBody: formatHTMLBody(message),
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Zoho-enczapikey "+m.c.Mailer.APIKey)

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			m.l.Error("mailer: failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("api returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (m *Mailer) sendSMTP(to, topic, message string) error {
	sc := m.c.Mailer.SMTP
	if strings.TrimSpace(sc.Host) == "" || sc.Port == 0 {
		return errors.New("smtp host or port is not configured")
	}

	addr := net.JoinHostPort(sc.Host, strconv.Itoa(sc.Port))
	tlsConfig := &tls.Config{ServerName: sc.Host}
	dialer := &net.Dialer{Timeout: m.timeout}

	var conn net.Conn
	var err error
	if sc.Encryption == encryptionTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	if err = conn.SetDeadline(time.Now().Add(m.timeout)); err != nil {
		_ = conn.Close()
		return err
	}

	client, err := smtp.NewClient(conn, sc.Host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer func() { _ = client.Close() }()

	if sc.Encryption != encryptionTLS {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("smtp server does not support STARTTLS")
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}

	if sc.Username != "" {
		if err = client.Auth(smtp.PlainAuth("", sc.Username, sc.Password, sc.Host)); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
	}

	if err = client.Mail(m.c.Mailer.Sender); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(m.buildMessage(to, topic, message)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func (m *Mailer) buildMessage(to, topic, message string) []byte {
	from := mail.Address{Name: config.AppName, Address: m.c.Mailer.Sender}
	domain := m.c.Mailer.Sender[strings.LastIndex(m.c.Mailer.Sender, "@")+1:]

	headers := []string{
		"From: " + from.String(),
		"To: " + (&mail.Address{Address: to}).String(),
		"Subject: " + mime.QEncoding.Encode("utf-8", topic),
		"Date: " + time.Now().Format(time.RFC1123Z),
		"Message-ID: <" + uuid.NewString() + "@" + domain + ">",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
	}

	body := strings.ReplaceAll(strings.ReplaceAll(message, "\r\n", "\n"), "\n", "\r\n")

	return []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body + "\r\n")
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

func (m *Mailer) SendOtp(to, otp string) {
	message := strings.Join([]string{
		"Dear user,",
		"Your One-Time Password (OTP) is:",
		"",
		otp,
		"",
		"Please enter this code within the next 3 minutes to complete your email verification process.",
		"For your security, do not share this code with anyone.",
		"If you did not request this code, please ignore this email.",
		"",
		"Thank you for using our service!",
		"https://nightell.neatplex.com",
	}, "\r\n")
	m.Send(to, "Nightell OTP (one-time password)", message)
}

func (m *Mailer) SendDeleteAccount(to, username, link string) {
	message := strings.Join([]string{
		"Dear " + username + ",",
		"We have received a request to delete your account associated with this email address.",
		"Please note that deleting your account is a permanent action and cannot be undone. " +
			"Get your data, including any saved preferences and history, will be permanently erased.",
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
	timeout := time.Duration(c.HttpClient.Timeout) * time.Millisecond
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Mailer{
		l:       l,
		c:       c,
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

func formatHTMLBody(message string) string {
	return message
}
