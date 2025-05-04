package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/peatch-io/peatch/internal/contract"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockPhotoUploader struct {
	UploadedFiles map[string]string // Stores key:content for verification
}

func (m *MockPhotoUploader) UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error {
	// Read the file content to simulate upload
	content, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if m.UploadedFiles == nil {
		m.UploadedFiles = make(map[string]string)
	}
	// Store the content for verification
	m.UploadedFiles[key] = string(content)
	return nil
}

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
	e := setupDependencies(t)

	authResp, err := authHelper(t, e, 927635965, "mkkksim", "Maksim")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	token := authResp.Token

	t.Run("UploadPNG", func(t *testing.T) {
		rec := performMultipartRequest(t, e, "photo", "image.png", "mock png content", token, http.StatusOK)
		var status contract.StatusResponse
		err := json.Unmarshal(rec.Body.Bytes(), &status)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		assert.True(t, status.Success, true)
	})

	t.Run("InvalidFileType", func(t *testing.T) {
		rec := performMultipartRequest(t, e, "photo", "test.txt", "text content", token, http.StatusBadRequest)
		errResp := parseResponse[contract.ErrorResponse](t, rec)
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
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected status 400")
		errResp := parseResponse[contract.ErrorResponse](t, rec)
		assert.Equal(t, "failed to get photo from form", errResp.Error)
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		rec := performMultipartRequest(t, e, "photo", "test.jpg", "mock jpg content", "", http.StatusUnauthorized)
		errResp := parseResponse[contract.ErrorResponse](t, rec)
		assert.Equal(t, handler.ErrAuthInvalid, errResp.Error)
	})
}
