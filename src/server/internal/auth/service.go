package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type Claims struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
	jwt.RegisteredClaims
}

type Service struct {
	secret   []byte
	mu       sync.RWMutex
	sessions map[string]time.Time
	now      func() time.Time
}

func NewService(secret string) *Service {
	return &Service{
		secret:   []byte(secret),
		sessions: make(map[string]time.Time),
		now:      time.Now,
	}
}

func (s *Service) Create(userID int) (string, error) {
	sessionID, err := randomID()
	if err != nil {
		return "", err
	}
	expiresAt := s.now().Add(15 * time.Minute)
	claims := Claims{
		SessionID: sessionID,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(s.now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	s.mu.Lock()
	s.sessions[sessionID] = expiresAt
	s.mu.Unlock()
	return signed, nil
}

func (s *Service) Validate(authorization string) (Claims, error) {
	tokenValue := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	if tokenValue == "" || tokenValue == authorization {
		return Claims{}, ErrInvalidToken
	}

	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	s.mu.RLock()
	expiresAt, exists := s.sessions[claims.SessionID]
	s.mu.RUnlock()
	if !exists || s.now().After(expiresAt) {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) Revoke(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
}

func randomID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("create session id: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
