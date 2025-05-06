package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/nanoid"
	"net/http"
	"path"
	"strings"
)

var allowedPhotoExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
	".svg":  {},
}

var ErrInvalidPhotoFormat = "Invalid photo format"

// handleUserAvatar godoc
// @Summary Upload user photo
// @Description Upload a photo for the authenticated user to S3 and store record in database
// @Tags photos
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Photo file to upload"
// @Success 200 {object} contract.StatusResponse
// @Failure 400 {object} contract.ErrorResponse
// @Failure 401 {object} contract.ErrorResponse
// @Failure 500 {object} contract.ErrorResponse
// @Router /api/users/avatar [post]
func (h *handler) handleUserAvatar(c echo.Context) error {
	userID := getUserID(c)

	file, err := c.FormFile("photo")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get photo from form").WithInternal(err)
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to open photo file").WithInternal(err)
	}
	defer src.Close()

	fileExtension := strings.ToLower(path.Ext(file.Filename))
	if _, ok := allowedPhotoExtensions[fileExtension]; !ok {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidPhotoFormat)
	}

	photoID := nanoid.Must()

	filename := fmt.Sprintf("photos/%s/%s%s", userID, photoID, fileExtension)

	if err = h.s3Client.UploadFile(context.TODO(), filename, src, file.Header.Get("Content-Type")); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to upload photo to S3").WithInternal(err)
	}

	err = h.storage.UpdateUserAvatarURL(c.Request().Context(), userID, filename)

	return c.JSON(http.StatusOK, contract.StatusResponse{Success: true})
}
