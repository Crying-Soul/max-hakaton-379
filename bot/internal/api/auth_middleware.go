package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"maxBot/internal/auth"
)

type contextKey string

const (
	ctxUserKey  contextKey = "max.auth.user"
	ctxTokenKey contextKey = "max.auth.token"
)

type authMiddleware struct {
	validator *auth.Validator
}

func newAuthMiddleware(validator *auth.Validator) *authMiddleware {
	if validator == nil {
		return nil
	}
	return &authMiddleware{validator: validator}
}

func (m *authMiddleware) requireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{Message: "требуется заголовок Authorization"})
			return
		}

		user, err := m.validator.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{Message: "некорректный или истёкший токен"})
			return
		}

		c.Set(string(ctxUserKey), user)
		c.Set(string(ctxTokenKey), token)
		c.Next()
	}
}

func extractBearerToken(header string) string {
	const prefix = "Bearer "
	trimmed := strings.TrimSpace(header)
	if trimmed == "" || !strings.HasPrefix(trimmed, prefix) {
		return ""
	}
	return strings.TrimSpace(trimmed[len(prefix):])
}

func getAuthenticatedUser(c *gin.Context) (auth.MaxUser, bool) {
	if value, ok := c.Get(string(ctxUserKey)); ok {
		if user, ok := value.(auth.MaxUser); ok {
			return user, true
		}
	}
	return auth.MaxUser{}, false
}

func getTokenFromContext(c *gin.Context) (string, bool) {
	if value, ok := c.Get(string(ctxTokenKey)); ok {
		if token, ok := value.(string); ok {
			return token, true
		}
	}
	return "", false
}
