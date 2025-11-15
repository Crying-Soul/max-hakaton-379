package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"maxBot/internal/auth"
)

func TestAuthMiddlewareAllowsValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator, secret := newTestValidator(t)
	token := signTestToken(t, secret, auth.MaxUser{ID: 99, FirstName: "Auth"})

	router := gin.New()
	router.GET("/secure", newAuthMiddleware(validator).requireUser(), func(c *gin.Context) {
		user, ok := getAuthenticatedUser(c)
		if !ok {
			t.Fatalf("user missing in context")
		}
		c.JSON(http.StatusOK, gin.H{"userID": user.ID})
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator, _ := newTestValidator(t)

	router := gin.New()
	router.GET("/secure", newAuthMiddleware(validator).requireUser(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func newTestValidator(t *testing.T) (*auth.Validator, string) {
	t.Helper()
	secret := "test-jwt-secret"
	validator, err := auth.NewValidator(auth.Config{
		BotToken:   "bot-token",
		JWTSecret:  secret,
		MaxAge:     time.Hour,
		SessionTTL: time.Hour,
	})
	if err != nil {
		t.Fatalf("init validator: %v", err)
	}
	return validator, secret
}

func signTestToken(t *testing.T, secret string, user auth.MaxUser) string {
	t.Helper()
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  fmt.Sprint(user.ID),
		"user": user,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}
