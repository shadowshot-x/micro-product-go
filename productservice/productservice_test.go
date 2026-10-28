package productservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestProductHandlers_MissingHeaders(t *testing.T) {
	ctrl := NewProductController(zap.NewNop())

	t.Run("GetAllProductByIdHandler_MissingId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product", nil)
		rec := httptest.NewRecorder()

		ctrl.GetAllProductByIdHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		if rec.Body.String() != "Id Missing" {
			t.Errorf("expected body %q, got %q", "Id Missing", rec.Body.String())
		}
	})

	t.Run("AddProductHandler_MissingName", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/product", nil)
		rec := httptest.NewRecorder()

		ctrl.AddProductHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		if rec.Body.String() != "name Missing" {
			t.Errorf("expected body %q, got %q", "name Missing", rec.Body.String())
		}
	})

	t.Run("DeleteProductHandler_MissingId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/product", nil)
		rec := httptest.NewRecorder()

		ctrl.DeleteProductHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		if rec.Body.String() != "Id Missing" {
			t.Errorf("expected body %q, got %q", "Id Missing", rec.Body.String())
		}
	})
}
