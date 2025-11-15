package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"maxBot/internal/model"
	"maxBot/internal/service"
)

type mapHandler struct {
	events service.EventService
}

func newMapHandler(events service.EventService) *mapHandler {
	return &mapHandler{events: events}
}

type mapEventsResponse struct {
	Data []model.MapEvent `json:"data"`
	Meta map[string]any   `json:"meta"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func (h *mapHandler) register(r *gin.RouterGroup, authMW *authMiddleware) {
	r.GET("/map/events", h.listEvents)
	usersGroup := r.Group("/map")
	if authMW != nil {
		usersGroup.Use(authMW.requireUser())
	}
	usersGroup.GET("/users/:userID/events", h.listUserEvents)
}

func (h *mapHandler) listEvents(c *gin.Context) {
	params, err := parseMapQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: err.Error()})
		return
	}

	items, err := h.events.ListEventsForMap(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "не удалось получить события"})
		return
	}

	resp := mapEventsResponse{
		Data: items,
		Meta: map[string]any{
			"limit":      params.Limit,
			"offset":     params.Offset,
			"count":      len(items),
			"radiusKm":   params.RadiusKm,
			"categories": params.CategoryIDs,
			"lat":        params.Lat,
			"lon":        params.Lon,
		},
	}

	c.JSON(http.StatusOK, resp)
}

func (h *mapHandler) listUserEvents(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, errorResponse{Message: "нельзя запрашивать чужие события"})
		return
	}

	params, err := parseMapQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: err.Error()})
		return
	}

	items, err := h.events.ListEventsForMapByVolunteer(c.Request.Context(), service.ListMapEventsForVolunteerParams{
		ListMapEventsParams: params,
		VolunteerID:         userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "не удалось получить события пользователя"})
		return
	}

	resp := mapEventsResponse{
		Data: items,
		Meta: map[string]any{
			"limit":      params.Limit,
			"offset":     params.Offset,
			"count":      len(items),
			"radiusKm":   params.RadiusKm,
			"categories": params.CategoryIDs,
			"lat":        params.Lat,
			"lon":        params.Lon,
			"userId":     userID,
		},
	}

	c.JSON(http.StatusOK, resp)
}

func parseFloatQuery(val string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(val), 64)
}

func parseCategoryIDs(c *gin.Context) ([]int32, error) {
	var raw []string
	if list := c.QueryArray("category_id"); len(list) > 0 {
		raw = append(raw, list...)
	}
	if list := c.Query("categories"); list != "" {
		raw = append(raw, strings.Split(list, ",")...)
	}
	if len(raw) == 0 {
		return nil, nil
	}
	result := make([]int32, 0, len(raw))
	for _, value := range raw {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		n, err := strconv.ParseInt(trimmed, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("category_id '%s' не число", trimmed)
		}
		result = append(result, int32(n))
	}
	return deduplicate(result), nil
}

func parseMapQueryParams(c *gin.Context) (service.ListMapEventsParams, error) {
	latStr := c.Query("lat")
	lat, err := parseFloatQuery(latStr)
	if err != nil {
		return service.ListMapEventsParams{}, fmt.Errorf("lat обязателен и должен быть числом")
	}
	lonStr := c.Query("lon")
	lon, err := parseFloatQuery(lonStr)
	if err != nil {
		return service.ListMapEventsParams{}, fmt.Errorf("lon обязателен и должен быть числом")
	}
	radiusStr := c.Query("radius_km")
	if strings.TrimSpace(radiusStr) == "" {
		return service.ListMapEventsParams{}, fmt.Errorf("radius_km обязателен")
	}
	radius, err := parseFloatQuery(radiusStr)
	if err != nil || radius <= 0 {
		return service.ListMapEventsParams{}, fmt.Errorf("radius_km обязателен и должен быть > 0")
	}

	limit := parseInt32Bound(c.Query("limit"), 1, 100, 50)
	offset := parseInt32Bound(c.Query("offset"), 0, 10_000, 0)
	categories, err := parseCategoryIDs(c)
	if err != nil {
		return service.ListMapEventsParams{}, err
	}

	return service.ListMapEventsParams{
		Offset:      offset,
		Limit:       limit,
		Lat:         lat,
		Lon:         lon,
		RadiusKm:    radius,
		CategoryIDs: categories,
	}, nil
}

func deduplicate(items []int32) []int32 {
	if len(items) == 0 {
		return nil
	}
	sort.Slice(items, func(i, j int) bool { return items[i] < items[j] })
	result := make([]int32, 0, len(items))
	prev := items[0]
	result = append(result, prev)
	for i := 1; i < len(items); i++ {
		if items[i] == prev {
			continue
		}
		prev = items[i]
		result = append(result, prev)
	}
	return result
}
