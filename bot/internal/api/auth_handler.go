package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"maxBot/internal/auth"
)

type authHandler struct {
	validator *auth.Validator
}

func newAuthHandler(validator *auth.Validator) *authHandler {
	return &authHandler{validator: validator}
}

func (h *authHandler) register(r *gin.RouterGroup) {
	r.GET("/auth/session", h.createSession)
}

type sessionResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expiresIn"`
	User      auth.MaxUser `json:"user"`
}

func (h *authHandler) createSession(c *gin.Context) {
	initData := strings.TrimSpace(c.Query("session"))
	if initData == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "initData обязателен"})
		return
	}

	result, err := h.validator.ValidateAndIssue(initData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessionResponse{
		Token:     result.Token,
		ExpiresIn: int64(h.validator.SessionTTL().Seconds()),
		User:      result.User,
	})
}
