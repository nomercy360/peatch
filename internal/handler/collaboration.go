package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/nanoid"
	"net/http"
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
// @Success 200 {array} contract.CollaborationResponse
// @Router /api/collaborations [get]
func (h *handler) handleListCollaborations(c echo.Context) error {
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)
	search := c.QueryParam("search")
	uid := getUserID(c)

	query := db.CollaborationQuery{
		Page:     page,
		Limit:    limit,
		Search:   search,
		ViewerID: uid,
	}

	collaborations, err := h.storage.ListCollaborations(c.Request().Context(), query)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaborations").WithInternal(err)
	}

	resp := make([]contract.CollaborationResponse, len(collaborations))
	for i, collaboration := range collaborations {
		resp[i] = contract.ToCollaborationResponse(collaboration)
	}

	return c.JSON(http.StatusOK, resp)
}

// handleGetCollaboration godoc
// @Summary Get collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param id path int true "Collaboration ID"
// @Success 200 {object} contract.CollaborationResponse
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

	return c.JSON(http.StatusOK, contract.ToCollaborationResponse(collaboration))
}

// handleCreateCollaboration godoc
// @Summary Create collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param collaboration body contract.CreateCollaboration true "Collaboration data"
// @Success 201 {object} contract.CollaborationResponse
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

	res, err := h.storage.GetCollaborationByID(c.Request().Context(), uid, collaboration.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collaboration not found")
	}

	return c.JSON(http.StatusCreated, contract.ToCollaborationResponse(res))
}

// handleUpdateCollaboration godoc
// @Summary Update collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param collaboration body contract.CreateCollaboration true "Collaboration data"
// @Success 200 {object} contract.CollaborationResponse
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

	collaboration, err := h.storage.GetCollaborationByID(c.Request().Context(), uid, cid)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collaboration not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collaboration").WithInternal(err)
	}

	return c.JSON(http.StatusOK, contract.ToCollaborationResponse(collaboration))
}
