package otp

import (
	"github.com/neatplex/nightell-core/internal/mailer"
	"math/rand"
	"strconv"
	"time"
)

const min3 = time.Duration(2) * time.Minute

type Service struct {
	mailer    *mailer.Mailer
	passwords map[string]*password
}

type password struct {
	Value    string
	ExpireAt time.Time
}

func (s *Service) Email(email string) int {
	if p, found := s.passwords[email]; found {
		if p.ExpireAt.After(time.Now()) {
			return s.ttl(p.ExpireAt)
		}
	}

	s.passwords[email] = &password{
		Value:    strconv.Itoa(rand.Intn(899999) + 100000),
		ExpireAt: time.Now().Add(min3),
	}

	s.mailer.SendOtp(email, s.passwords[email].Value)

	return s.ttl(s.passwords[email].ExpireAt)
}

func (s *Service) Check(identifier string, password string) bool {
	if p, found := s.passwords[identifier]; found {
		if p.ExpireAt.Before(time.Now().Add(min3)) {
			if p.Value == password {
				delete(s.passwords, identifier)
				return true
			}
		} else {
			delete(s.passwords, identifier)
		}
	}
	return false
}

func (s *Service) ttl(t time.Time) int {
	return int(time.Until(t).Seconds())
}

func New(m *mailer.Mailer) *Service {
	return &Service{
		passwords: make(map[string]*password),
		mailer:    m,
	}
}
