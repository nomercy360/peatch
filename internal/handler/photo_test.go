package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/testutils"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performMultipartRequest(t *testing.T, e *echo.Echo, fieldName, fileName, fileContent, token string, expectedStatus int) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.WriteString(part, fileContent)
	if err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/users/avatar", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d, body: %s", expectedStatus, rec.Code, rec.Body.String())
	}
	return rec
}

func TestUploadPhoto(t *testing.T) {
	ts := testutils.SetupTestEnvironment(t)
	defer ts.Teardown()

	authResp, err := testutils.AuthHelper(t, ts.Echo, 927635965, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	token := authResp.Token

	t.Run("UploadPNG", func(t *testing.T) {
		rec := performMultipartRequest(t, ts.Echo, "photo", "image.png", "mock png content", token, http.StatusOK)
		var status contract.StatusResponse
		err := json.Unmarshal(rec.Body.Bytes(), &status)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		assert.True(t, status.Success, true)
	})

	t.Run("InvalidFileType", func(t *testing.T) {
		rec := performMultipartRequest(t, ts.Echo, "photo", "test.txt", "text content", token, http.StatusBadRequest)
		errResp := testutils.ParseResponse[contract.ErrorResponse](t, rec)
		assert.Equal(t, handler.ErrInvalidPhotoFormat, errResp.Error)
	})

	t.Run("MissingFile", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/users/avatar", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		rec := httptest.NewRecorder()
		ts.Echo.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status 400")
		errResp := testutils.ParseResponse[contract.ErrorResponse](t, rec)
		assert.Equal(t, "failed to get photo from form", errResp.Error)
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		rec := performMultipartRequest(t, ts.Echo, "photo", "test.jpg", "mock jpg content", "", http.StatusUnauthorized)
		// Just check the status code since middleware returns different error messages
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
