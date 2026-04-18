package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhegas/saas_gangsta/internal/modules/auth/dto"
	"github.com/gin-gonic/gin"
)

type mockAuthUsecase struct {
	registerRes *dto.RegisterResponse
	loginRes    *dto.LoginResponse
	refreshRes  *dto.LoginResponse
	meRes       *dto.MeResponse
	registerErr error
	loginErr    error
	refreshErr  error
	logoutErr   error
	meErr       error
}

func (m *mockAuthUsecase) Register(_ context.Context, _ dto.RegisterRequest) (*dto.RegisterResponse, error) {
	return m.registerRes, m.registerErr
}

func (m *mockAuthUsecase) Login(_ context.Context, _ dto.LoginRequest) (*dto.LoginResponse, error) {
	return m.loginRes, m.loginErr
}

func (m *mockAuthUsecase) Refresh(_ context.Context, _ dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	return m.refreshRes, m.refreshErr
}

func (m *mockAuthUsecase) Logout(_ context.Context, _ string, _ dto.LogoutRequest) error {
	return m.logoutErr
}

func (m *mockAuthUsecase) Me(_ context.Context, _ string) (*dto.MeResponse, error) {
	return m.meRes, m.meErr
}

func TestLoginHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{loginRes: &dto.LoginResponse{AccessToken: "a", RefreshToken: "r"}})
	r.POST("/login", h.Login)

	payload, _ := json.Marshal(dto.LoginRequest{Email: "user@test.local", Password: "secret123"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRegisterHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{registerRes: &dto.RegisterResponse{User: dto.UserResponse{ID: "u-1"}}})
	r.POST("/register", h.Register)

	payload, _ := json.Marshal(dto.RegisterRequest{Email: "user@test.local", Password: "secret123", FullName: "User Test"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestRefreshHandlerValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{})
	r.POST("/refresh", h.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestMeHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{meRes: &dto.MeResponse{User: dto.UserResponse{ID: "u-1"}}})
	r.GET("/me", func(c *gin.Context) {
		c.Set("userId", "u-1")
		h.Me(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRefreshHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{refreshRes: &dto.LoginResponse{AccessToken: "a", RefreshToken: "r"}})
	r.POST("/refresh", h.Refresh)

	payload, _ := json.Marshal(dto.RefreshTokenRequest{RefreshToken: "token"})
	req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestLogoutHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{})
	r.POST("/logout", func(c *gin.Context) {
		c.Set("userId", "u-1")
		h.Logout(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBufferString(`{"refreshToken":"r"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestMeHandlerUnauthorizedWhenUsecaseFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewAuthHandler(&mockAuthUsecase{meErr: errors.New("unauthorized")})
	r.GET("/me", func(c *gin.Context) {
		h.Me(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for non-app error mapping, got %d", w.Code)
	}
}
