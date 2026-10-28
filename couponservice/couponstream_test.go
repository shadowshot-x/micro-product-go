package couponservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestCouponStream_MissingHeaders(t *testing.T) {
	ctrl := NewCouponStreamController(zap.NewNop(), nil)

	t.Run("AddCouponList_MissingName", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/coupon/addcoupon", nil)
		rec := httptest.NewRecorder()

		ctrl.AddCouponList(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		if rec.Body.String() != "name Missing" {
			t.Errorf("expected body %q, got %q", "name Missing", rec.Body.String())
		}
	})

	t.Run("PurgeStream_MissingRegion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/coupon/purge", nil)
		rec := httptest.NewRecorder()

		ctrl.PurgeStream(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
		if rec.Body.String() != "region Missing" {
			t.Errorf("expected body %q, got %q", "region Missing", rec.Body.String())
		}
	})
}
