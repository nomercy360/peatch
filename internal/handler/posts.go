package handler

import (
	"github.com/labstack/echo/v4"
	svc "github.com/peatch-io/peatch/internal/service"
	"net/http"
	"strconv"
)

func (h *handler) handleCreatePost(c echo.Context) error {
	var req svc.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userID := getUserID(c)

	res, err := h.svc.CreatePost(userID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

// handleGetPost godoc
// @Summary Find post by id
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} Post
// @Router /api/posts/{id} [get]
func (h *handler) handleGetPost(c echo.Context) error {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	uid := getUserID(c)

	res, err := h.svc.GetPostByID(uid, postID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) handleUpdatePost(c echo.Context) error {
	var req svc.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userID := getUserID(c)
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	res, err := h.svc.UpdatePost(userID, postID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
