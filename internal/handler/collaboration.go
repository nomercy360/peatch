package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/nanoid"
	"net/http"
	"strconv"
	"time"
)

// handleListCollaborations godoc
// @Summary List collaborations
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param order query string false "Order by"
// @Success 200 {array} Collaboration
// @Router /api/collaborations [get]
func (h *handler) handleListCollaborations(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	uid := getUserID(c)

	query := db.CollaborationQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
		UserID: uid,
	}

	collaborations, err := h.storage.ListCollaborations(c.Request().Context(), query)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaborations").WithInternal(err)
	}

	return c.JSON(http.StatusOK, collaborations)
}

// handleGetCollaboration godoc
// @Summary Get collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param id path int true "Collaboration ID"
// @Success 200 {object} Collaboration
// @Router /api/collaborations/{id} [get]
func (h *handler) handleGetCollaboration(c echo.Context) error {
	id := c.Param("id")
	uid := getUserID(c)

	collaboration, err := h.storage.GetCollaborationByID(c.Request().Context(), uid, id)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collaboration not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaboration").WithInternal(err)
	}

	return c.JSON(http.StatusOK, collaboration)
}

// handleCreateCollaboration godoc
// @Summary Create collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param collaboration body CreateCollaboration true "Collaboration data"
// @Success 201 {object} Collaboration
// @Router /api/collaborations [post]
func (h *handler) handleCreateCollaboration(c echo.Context) error {
	var req contract.CreateCollaboration
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	uid := getUserID(c)
	now := time.Now()

	collaboration := db.Collaboration{
		ID:          nanoid.Must(),
		UserID:      uid,
		Title:       req.Title,
		Description: req.Description,
		IsPayable:   req.IsPayable,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.storage.CreateCollaboration(
		c.Request().Context(),
		collaboration,
		req.BadgeIDs,
		req.OpportunityID,
		req.LocationID,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "create failed").WithInternal(err)
	}

	return c.JSON(http.StatusCreated, collaboration)
}

// handleUpdateCollaboration godoc
// @Summary Update collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param collaboration body Collaboration true "Collaboration data"
// @Success 200 {object} Collaboration
// @Router /api/collaborations/{id} [put]
func (h *handler) handleUpdateCollaboration(c echo.Context) error {
	cid := c.Param("id")
	uid := getUserID(c)

	var req contract.CreateCollaboration
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidRequest).WithInternal(err)
	}

	collab := db.Collaboration{
		ID:          cid,
		UserID:      uid,
		Title:       req.Title,
		Description: req.Description,
		IsPayable:   req.IsPayable,
	}

	if err := h.storage.UpdateCollaboration(
		c.Request().Context(),
		collab,
		req.BadgeIDs,
		req.OpportunityID,
		req.LocationID,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "update failed").WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
