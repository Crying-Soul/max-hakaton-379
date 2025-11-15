package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"maxBot/internal/service"
)

type userHandler struct {
	users service.UserService
}

func newUserHandler(users service.UserService) *userHandler {
	if users == nil {
		return nil
	}
	return &userHandler{users: users}
}

type userLocationResponse struct {
	UserID    int64    `json:"userId"`
	Lat       *float64 `json:"lat"`
	Lon       *float64 `json:"lon"`
	UpdatedAt string   `json:"updatedAt"`
}

func (h *userHandler) register(r *gin.RouterGroup, authMW *authMiddleware) {
	if h == nil || authMW == nil {
		return
	}
	group := r.Group("/users")
	group.Use(authMW.requireUser())
	group.GET("/:userID/location", h.getUserLocation)
}

func (h *userHandler) getUserLocation(c *gin.Context) {
	userIDStr := strings.TrimSpace(c.Param("userID"))
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "userID обязателен"})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "userID должен быть положительным числом"})
		return
	}

	authUser, ok := getAuthenticatedUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Message: "требуется авторизация"})
		return
	}
	if authUser.ID != userID {
		c.JSON(http.StatusForbidden, errorResponse{Message: "нельзя запрашивать геолокацию другого пользователя"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.users.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse{Message: "пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "не удалось получить пользователя"})
		return
	}

	if user.LocationLat == nil || user.LocationLon == nil {
		c.JSON(http.StatusNotFound, errorResponse{Message: "у пользователя нет сохранённой геолокации"})
		return
	}

	resp := userLocationResponse{
		UserID:    user.ID,
		Lat:       user.LocationLat,
		Lon:       user.LocationLon,
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, resp)
}
