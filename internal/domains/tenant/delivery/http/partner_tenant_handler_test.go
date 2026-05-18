package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/gin-gonic/gin"
)

type mockPartnerTenantUsecase struct {
	createTenantRes *dto.CreatePartnerTenantResponse
	listTenantsRes  *dto.ListPartnerTenantsResponse
	createTenantErr error
	listTenantsErr  error
}

func (m *mockPartnerTenantUsecase) CreatePartnerTenant(_ context.Context, _ string, _ dto.CreatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error) {
	return m.createTenantRes, m.createTenantErr
}

func (m *mockPartnerTenantUsecase) ListPartnerTenants(_ context.Context, _ string) (*dto.ListPartnerTenantsResponse, error) {
	return m.listTenantsRes, m.listTenantsErr
}

func (m *mockPartnerTenantUsecase) SoftDeletePartnerTenant(_ context.Context, _ string, _ string) error {
	return nil
}

func TestCreatePartnerTenantHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPartnerTenantHandler(&mockPartnerTenantUsecase{createTenantRes: &dto.CreatePartnerTenantResponse{Tenant: dto.PartnerTenantResponse{ID: "t-1", Name: "Warung A", Slug: "warung-a", Status: "active", IsOwner: true}}})
	r.POST("/partner/tenants", func(c *gin.Context) {
		c.Set("userId", "u-1")
		h.CreatePartnerTenant(c)
	})

	payload, _ := json.Marshal(dto.CreatePartnerTenantRequest{Name: "Warung A"})
	req := httptest.NewRequest(http.MethodPost, "/partner/tenants", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestListPartnerTenantsHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewPartnerTenantHandler(&mockPartnerTenantUsecase{listTenantsRes: &dto.ListPartnerTenantsResponse{Tenants: []dto.PartnerTenantResponse{{ID: "t-1", Name: "Warung A", Slug: "warung-a", Status: "active", IsOwner: true}}}})
	r.GET("/partner/tenants", func(c *gin.Context) {
		c.Set("userId", "u-1")
		h.ListPartnerTenants(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/partner/tenants", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
