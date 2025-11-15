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
	r.POST("/auth/session", h.createSession)
}

type sessionRequest struct {
	InitData string `json:"initData"`
}

type sessionResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expiresIn"`
	User      auth.MaxUser `json:"user"`
}

func (h *authHandler) createSession(c *gin.Context) {
	var req sessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "initData обязателен"})
		return
	}

	initData := strings.TrimSpace(req.InitData)
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
