package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"net/http"
	"strconv"
)

func parseIntQuery(c echo.Context, key string, defaultValue int) int {
	value, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || value < 0 {
		return defaultValue
	}
	return value
}

// handleSearchLocations godoc
// @Summary List cities
// @Tags cities
// @Accept  json
// @Produce  json
// @Success 200 {array} contract.CityResponse
// @Router /api/locations [get]
func (h *handler) handleSearchLocations(c echo.Context) error {
	query := c.QueryParam("search")

	limit := parseIntQuery(c, "limit", 10)
	page := parseIntQuery(c, "page", 1)

	skip := (page - 1) * limit

	cities, err := h.storage.SearchCities(c.Request().Context(), query, limit, skip)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to search locations").WithInternal(err)
	}

	resp := make([]contract.CityResponse, len(cities))
	for i, city := range cities {
		resp[i] = contract.ToCityResponse(city)
	}

	return c.JSON(http.StatusOK, cities)
}
