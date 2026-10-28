package clientclaims

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestClientClaims_Handlers(t *testing.T) {
	logger := zap.NewNop()

	t.Run("UploadFile_InvalidMultipart", func(t *testing.T) {
		ctrl := NewUploadController(logger)
		req := httptest.NewRequest(http.MethodPost, "/claims/upload", strings.NewReader("not a multipart body"))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")
		rec := httptest.NewRecorder()

		ctrl.UploadFile(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("DownloadFile_MissingEmail", func(t *testing.T) {
		ctrl := NewDownloadController(logger)
		req := httptest.NewRequest(http.MethodGet, "/claims/download", nil)
		rec := httptest.NewRecorder()

		ctrl.DownloadFile(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		if rec.Body.String() != "Email Missing" {
			t.Errorf("expected 'Email Missing', got %q", rec.Body.String())
		}
	})
}
