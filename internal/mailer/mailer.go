package mailer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"go.uber.org/zap"
)

const (
	defaultMailerEndpoint = "https://api.zeptomail.eu/v1.1/email"
)

type Mailer struct {
	c      *config.Config
	l      *logger.Logger
	client *http.Client
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
	if m.client == nil {
		m.l.Error("mailer: http client not initialized")
		return
	}

	if strings.TrimSpace(m.c.Mailer.Sender) == "" {
		m.l.Error("mailer: sender email is not configured")
		return
	}

	if strings.TrimSpace(m.c.Mailer.APIKey) == "" {
		m.l.Error("mailer: api key is not configured")
		return
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
		m.l.Error("mailer: failed to marshal request", zap.Error(err))
		return
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		m.l.Error("mailer: failed to create request", zap.Error(err))
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Zoho-enczapikey "+m.c.Mailer.APIKey)

	resp, err := m.client.Do(req)
	if err != nil {
		m.l.Info("mailer: failed", zap.String("to", to), zap.Error(err))
		return
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			m.l.Error("mailer: failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		m.l.Error("mailer: api returned error",
			zap.String("to", to),
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(body)))
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
		l:      l,
		c:      c,
		client: &http.Client{Timeout: timeout},
	}
}

func formatHTMLBody(message string) string {
	return message
}
