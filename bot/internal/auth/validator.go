package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"maxBot/internal/auth/initdatahash"
)

// Config controls validation and JWT generation.
type Config struct {
	BotToken   string
	JWTSecret  string
	MaxAge     time.Duration
	SessionTTL time.Duration
}

// Validator validates MAX init data and issues JWT tokens.
type Validator struct {
	botToken   string
	jwtSecret  []byte
	maxAge     time.Duration
	sessionTTL time.Duration
	now        func() time.Time
}

type maxClaims struct {
	User MaxUser `json:"user"`
	jwt.RegisteredClaims
}

// NewValidator constructs a Validator ensuring config defaults.
func NewValidator(cfg Config) (*Validator, error) {
	if cfg.BotToken == "" {
		return nil, errors.New("bot token is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if cfg.MaxAge <= 0 {
		cfg.MaxAge = time.Hour
	}
	if cfg.SessionTTL <= 0 {
		cfg.SessionTTL = 24 * time.Hour
	}
	return &Validator{
		botToken:   cfg.BotToken,
		jwtSecret:  []byte(cfg.JWTSecret),
		maxAge:     cfg.MaxAge,
		sessionTTL: cfg.SessionTTL,
		now:        time.Now,
	}, nil
}

// AuthResult contains parsed payload and issued JWT token.
type AuthResult struct {
	Token     string
	User      MaxUser
	RawParams url.Values
}

// MaxUser mirrors minimal data extracted from initData.
type MaxUser struct {
	ID           int64   `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Username     *string `json:"username"`
	LanguageCode *string `json:"language_code"`
	PhotoURL     *string `json:"photo_url"`
}

// ValidateAndIssue validates MAX init data and returns a signed JWT.
func (v *Validator) ValidateAndIssue(rawInitData string) (AuthResult, error) {
	decoded, err := url.QueryUnescape(rawInitData)
	if err != nil {
		return AuthResult{}, fmt.Errorf("init data decode: %w", err)
	}

	values, err := url.ParseQuery(decoded)
	if err != nil {
		return AuthResult{}, fmt.Errorf("init data parse: %w", err)
	}

	if err := v.verifyHash(values, values.Get("hash")); err != nil {
		return AuthResult{}, err
	}

	if err := v.verifyFreshness(values.Get("auth_date")); err != nil {
		return AuthResult{}, err
	}

	user, err := parseUser(values.Get("user"))
	if err != nil {
		return AuthResult{}, err
	}

	token, err := v.issueJWT(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, User: user, RawParams: values}, nil
}

func (v *Validator) verifyHash(values url.Values, provided string) error {
	normalized := initdatahash.NormalizeHash(provided)
	if normalized == "" {
		return errors.New("hash parameter missing")
	}

	expected, _, err := initdatahash.Compute(values, v.botToken)
	if err != nil {
		return err
	}

	if !initdatahash.HashesEqual(normalized, expected) {
		return errors.New("init data signature mismatch")
	}
	return nil
}

func (v *Validator) verifyFreshness(authDate string) error {
	return initdatahash.CheckFreshness(authDate, v.maxAge, v.now())
}

func (v *Validator) issueJWT(user MaxUser) (string, error) {
	now := v.now()
	claims := &maxClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(v.sessionTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(v.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return signed, nil
}

// ParseToken validates a JWT issued by this validator and extracts the user payload.
func (v *Validator) ParseToken(token string) (MaxUser, error) {
	claims := &maxClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %T", t.Method)
		}
		return v.jwtSecret, nil
	})
	if err != nil {
		return MaxUser{}, fmt.Errorf("parse token: %w", err)
	}
	if !parsed.Valid {
		return MaxUser{}, errors.New("token invalid")
	}
	return claims.User, nil
}

// SessionTTL returns configured session lifetime.
func (v *Validator) SessionTTL() time.Duration {
	return v.sessionTTL
}

func parseUser(raw string) (MaxUser, error) {
	if raw == "" {
		return MaxUser{}, errors.New("user payload missing")
	}
	var user MaxUser
	if err := json.Unmarshal([]byte(raw), &user); err != nil {
		return MaxUser{}, fmt.Errorf("user decode: %w", err)
	}
	return user, nil
}

func buildDataCheckString(values url.Values) string {
	return initdatahash.BuildDataCheckString(values)
}
