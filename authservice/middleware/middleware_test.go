package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/shadowshot-x/micro-product-go/authservice/jwt"
	"go.uber.org/zap"
)

func TestTokenValidationMiddleware(t *testing.T) {
	logger := zap.NewNop()
	middleware := NewTokenMiddleware(logger)

	nextHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("next handler reached"))
	})

	handler := middleware.TokenValidationMiddleware(nextHandler)

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
		if rec.Body.String() != "Token Missing" {
			t.Errorf("expected body %q, got %q", "Token Missing", rec.Body.String())
		}
	})

	t.Run("Empty Secret", func(t *testing.T) {
		origSecret := os.Getenv("JWT_SECRET")
		os.Unsetenv("JWT_SECRET")
		defer os.Setenv("JWT_SECRET", origSecret)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Token", "some.token.here")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
		if rec.Body.String() != "Internal Server Error" {
			t.Errorf("expected body %q, got %q", "Internal Server Error", rec.Body.String())
		}
	})

	t.Run("Corrupt Token", func(t *testing.T) {
		origSecret := os.Getenv("JWT_SECRET")
		os.Setenv("JWT_SECRET", "testsecret")
		defer os.Setenv("JWT_SECRET", origSecret)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Token", "not-a-valid-token")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
		if rec.Body.String() != jwt.CORRUPT_TOKEN {
			t.Errorf("expected body %q, got %q", jwt.CORRUPT_TOKEN, rec.Body.String())
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		origSecret := os.Getenv("JWT_SECRET")
		os.Setenv("JWT_SECRET", "testsecret")
		defer os.Setenv("JWT_SECRET", origSecret)

		validClaims := jwt.ClaimsMap{
			Aud: "test",
			Iss: "test",
			Exp: fmt.Sprint(time.Now().Add(time.Hour).Unix()),
		}
		token, err := jwt.GenerateToken("HS256", validClaims, "different-secret")
		if err != nil {
			t.Fatalf("token generation failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Token", token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
		if rec.Body.String() != jwt.INVALID_TOKEN {
			t.Errorf("expected body %q, got %q", jwt.INVALID_TOKEN, rec.Body.String())
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		origSecret := os.Getenv("JWT_SECRET")
		os.Setenv("JWT_SECRET", "testsecret")
		defer os.Setenv("JWT_SECRET", origSecret)

		expiredClaims := jwt.ClaimsMap{
			Aud: "test",
			Iss: "test",
			Exp: fmt.Sprint(time.Now().Add(-1 * time.Hour).Unix()),
		}
		token, err := jwt.GenerateToken("HS256", expiredClaims, "testsecret")
		if err != nil {
			t.Fatalf("token generation failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Token", token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
		if rec.Body.String() != jwt.EXPIRED_TOKEN {
			t.Errorf("expected body %q, got %q", jwt.EXPIRED_TOKEN, rec.Body.String())
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		origSecret := os.Getenv("JWT_SECRET")
		os.Setenv("JWT_SECRET", "testsecret")
		defer os.Setenv("JWT_SECRET", origSecret)

		validClaims := jwt.ClaimsMap{
			Aud: "test",
			Iss: "test",
			Exp: fmt.Sprint(time.Now().Add(time.Hour).Unix()),
		}
		token, err := jwt.GenerateToken("HS256", validClaims, "testsecret")
		if err != nil {
			t.Fatalf("token generation failed: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Token", token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
		if rec.Body.String() != "next handler reached" {
			t.Errorf("expected body %q, got %q", "next handler reached", rec.Body.String())
		}
	})
}
