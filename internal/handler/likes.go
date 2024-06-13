package handler

import (
	"github.com/labstack/echo/v4"
	svc "github.com/peatch-io/peatch/internal/service"
	"net/http"
)

func (h *handler) handleIncreaseLikeCount(c echo.Context) error {
	uid := getUserID(c)

	var likeRequest svc.LikeRequest

	if err := c.Bind(&likeRequest); err != nil {
		return err
	}

	if err := c.Validate(likeRequest); err != nil {
		return err
	}

	if err := h.svc.IncreaseLikeCount(uid, likeRequest); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) handleDecreaseLikeCount(c echo.Context) error {
	uid := getUserID(c)

	var likeRequest svc.LikeRequest

	if err := c.Bind(&likeRequest); err != nil {
		return err
	}

	if err := c.Validate(likeRequest); err != nil {
		return err
	}

	if err := h.svc.DecreaseLikeCount(uid, likeRequest); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
