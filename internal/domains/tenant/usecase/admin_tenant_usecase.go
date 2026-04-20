package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type adminTenantUsecase struct {
	repo domain.AdminTenantRepository
}

// newUUID menghasilkan UUID v4 menggunakan crypto/rand (tanpa dependency tambahan)
func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// NewAdminTenantUsecase adalah constructor untuk dependency injection
func NewAdminTenantUsecase(repo domain.AdminTenantRepository) domain.AdminTenantUsecase {
	return &adminTenantUsecase{repo: repo}
}

// GetAllTenants mengambil seluruh tenant aktif dari database
func (u *adminTenantUsecase) GetAllTenants(ctx context.Context) ([]dto.TenantResponse, error) {
	entities, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.TenantResponse
	for _, e := range entities {
		responses = append(responses, toTenantResponse(e))
	}

	return responses, nil
}

// GetTenantByID mengambil detail satu tenant berdasarkan ID
func (u *adminTenantUsecase) GetTenantByID(ctx context.Context, tenantID string) (*dto.TenantResponse, error) {
	entity, err := u.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	res := toTenantResponse(*entity)
	return &res, nil
}

// CreateTenant membuat tenant baru dengan default status "active"
func (u *adminTenantUsecase) CreateTenant(ctx context.Context, req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	// Default status jika tidak disertakan
	status := req.Status
	if status == "" {
		status = "active"
	}

	now := time.Now()
	entity := &domain.TenantEntity{
		ID:        newUUID(),
		Name:      req.Name,
		Slug:      req.Slug,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("gagal membuat tenant: %w", err)
	}

	res := toTenantResponse(*entity)
	return &res, nil
}

// UpdateTenant memperbarui data tenant (name, slug, status)
func (u *adminTenantUsecase) UpdateTenant(ctx context.Context, tenantID string, req dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	// Pastikan tenant ada terlebih dahulu
	existing, err := u.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant tidak ditemukan: %w", err)
	}

	// Terapkan perubahan hanya untuk field yang dikirim
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Slug != "" {
		existing.Slug = req.Slug
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("gagal memperbarui tenant: %w", err)
	}

	// Ambil data terbaru dari DB supaya response akurat
	updated, err := u.repo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	res := toTenantResponse(*updated)
	return &res, nil
}

// SoftDeleteTenant menandai tenant sebagai terhapus (mengisi deleted_at)
func (u *adminTenantUsecase) SoftDeleteTenant(ctx context.Context, tenantID string) error {
	// Pastikan tenant ada sebelum dihapus
	if _, err := u.repo.FindByID(ctx, tenantID); err != nil {
		return fmt.Errorf("tenant tidak ditemukan: %w", err)
	}

	return u.repo.SoftDelete(ctx, tenantID)
}

// UpdateTenantStatus memperbarui hanya status tenant
func (u *adminTenantUsecase) UpdateTenantStatus(ctx context.Context, tenantID string, status string) error {
	return u.repo.UpdateStatus(ctx, tenantID, status)
}

// toTenantResponse adalah helper mapping entity → DTO
func toTenantResponse(e domain.TenantEntity) dto.TenantResponse {
	return dto.TenantResponse{
		ID:        e.ID,
		Name:      e.Name,
		Slug:      e.Slug,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}
