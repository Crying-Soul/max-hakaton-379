package auth

import (
	"encoding/json"
	"net/url"
	"strconv"
	"testing"
	"time"

	"maxBot/internal/auth/initdatahash"
)

func TestValidateAndIssueAndParseToken(t *testing.T) {
	cfg := Config{
		BotToken:   "test-token",
		JWTSecret:  "jwt-secret",
		MaxAge:     time.Hour,
		SessionTTL: time.Hour,
	}
	validator, err := NewValidator(cfg)
	if err != nil {
		t.Fatalf("failed to init validator: %v", err)
	}

	user := MaxUser{ID: 321, FirstName: "Test", LastName: "User"}
	initData := buildInitData(t, cfg.BotToken, user, time.Now())

	res, err := validator.ValidateAndIssue(initData)
	if err != nil {
		t.Fatalf("ValidateAndIssue returned error: %v", err)
	}

	if res.Token == "" {
		t.Fatalf("expected token to be issued")
	}
	if res.User.ID != user.ID {
		t.Fatalf("unexpected user id: got %d want %d", res.User.ID, user.ID)
	}

	parsedUser, err := validator.ParseToken(res.Token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}
	if parsedUser.ID != user.ID {
		t.Fatalf("parsed user mismatch: got %d want %d", parsedUser.ID, user.ID)
	}
}

func TestParseTokenFailsForInvalidValue(t *testing.T) {
	validator, err := NewValidator(Config{BotToken: "token", JWTSecret: "secret"})
	if err != nil {
		t.Fatalf("init validator: %v", err)
	}

	if _, err := validator.ParseToken("invalid-token"); err == nil {
		t.Fatalf("expected error for invalid token")
	}
}

func buildInitData(t *testing.T, botToken string, user MaxUser, now time.Time) string {
	t.Helper()
	values := url.Values{}
	values.Set("auth_date", strconv.FormatInt(now.Unix(), 10))
	values.Set("query_id", "query-123")
	encodedUser, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}
	values.Set("user", string(encodedUser))

	expected, _, err := initdatahash.Compute(values, botToken)
	if err != nil {
		t.Fatalf("compute hash: %v", err)
	}
	values.Set("hash", expected)

	return values.Encode()
}
