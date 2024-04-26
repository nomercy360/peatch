package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/db"
	"net/http"
	"strconv"
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
	orderBy := c.QueryParam("order")

	query := db.CollaborationQuery{
		Page:    page,
		Limit:   limit,
		OrderBy: db.CollectionQueryOrder(orderBy),
	}

	collaborations, err := h.svc.ListCollaborations(query)

	if err != nil {
		return err
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
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	collaboration, err := h.svc.GetCollaborationByID(id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, collaboration)
}

// handleCreateCollaboration godoc
// @Summary Create collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param collaboration body Collaboration true "Collaboration data"
// @Success 201 {object} Collaboration
// @Router /api/collaborations [post]
func (h *handler) handleCreateCollaboration(c echo.Context) error {
	var collaboration db.Collaboration
	if err := c.Bind(&collaboration); err != nil {
		return err
	}

	createdCollaboration, err := h.svc.CreateCollaboration(collaboration)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdCollaboration)
}

// handleUpdateCollaboration godoc
// @Summary Update collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param collaboration body Collaboration true "Collaboration data"
// @Success 200 {object} Collaboration
// @Router /api/collaborations [put]
func (h *handler) handleUpdateCollaboration(c echo.Context) error {
	var collaboration db.Collaboration
	if err := c.Bind(&collaboration); err != nil {
		return err
	}

	updatedCollaboration, err := h.svc.UpdateCollaboration(collaboration)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedCollaboration)
}

// handlePublishCollaboration godoc
// @Summary Publish collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param id path int true "Collaboration ID"
// @Success 200
// @Router /api/collaborations/{id}/publish [put]
func (h *handler) handlePublishCollaboration(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.PublishCollaboration(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// handleHideCollaboration godoc
// @Summary Hide collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param id path int true "Collaboration ID"
// @Success 200
// @Router /api/collaborations/{id}/hide [put]
func (h *handler) handleHideCollaboration(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.HideCollaboration(id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// handleDeleteCollaboration godoc
// @Summary Delete collaboration
// @Tags collaborations
// @Accept  json
// @Produce  json
// @Param id path int true "Collaboration ID"
// @Success 204
// @Router /api/collaborations/{id} [delete]
func (h *handler) handleCreateCollaborationRequest(c echo.Context) error {
	var request db.CollaborationRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	createdRequest, err := h.svc.CreateCollaborationRequest(request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdRequest)
}
