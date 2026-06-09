package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
)

type partnerTenantUsecase struct {
	repo domain.PartnerTenantRepository
}

func NewPartnerTenantUsecase(repo domain.PartnerTenantRepository) domain.PartnerTenantUsecase {
	return &partnerTenantUsecase{repo: repo}
}

func (u *partnerTenantUsecase) CreatePartnerTenant(ctx context.Context, userID string, req dto.CreatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Nama tenant wajib diisi", http.StatusBadRequest)
	}

	tenant, err := u.repo.CreateTenantForPartner(ctx, domain.CreatePartnerTenantInput{
		UserID:      userID,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Address:     strings.TrimSpace(req.Address),
		PhoneNumber: strings.TrimSpace(req.PhoneNumber),
		OpenHours:   strings.TrimSpace(req.OpenHours),
		LogoURL:     req.LogoURL,
		BannerURL:   req.BannerURL,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound)
		case errors.Is(err, repository.ErrUserNotPartner):
			return nil, apperrors.New("FORBIDDEN", "Hanya PARTNER yang dapat membuat tenant", http.StatusForbidden)
		default:
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat tenant partner", http.StatusInternalServerError)
		}
	}

	return &dto.CreatePartnerTenantResponse{
		Tenant: dto.PartnerTenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			Status:      tenant.Status,
			Description: tenant.Description,
			Address:     tenant.Address,
			PhoneNumber: tenant.PhoneNumber,
			OpenHours:   tenant.OpenHours,
			LogoURL:     tenant.LogoURL,
			BannerURL:   tenant.BannerURL,
			IsOwner:     tenant.IsOwner,
		},
	}, nil
}

func (u *partnerTenantUsecase) ListPartnerTenants(ctx context.Context, userID string) (*dto.ListPartnerTenantsResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized)
	}

	partner, err := u.repo.FindPartnerByID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data user", http.StatusInternalServerError)
	}
	if partner == nil || !partner.IsActive {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak ditemukan", http.StatusUnauthorized)
	}
	if partner.Role != "PARTNER" {
		return nil, apperrors.New("FORBIDDEN", "Hanya PARTNER yang dapat melihat tenant", http.StatusForbidden)
	}

	tenants, err := u.repo.ListTenantsByPartner(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar tenant partner", http.StatusInternalServerError)
	}

	items := make([]dto.PartnerTenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		items = append(items, dto.PartnerTenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			Status:      tenant.Status,
			Description: tenant.Description,
			Address:     tenant.Address,
			PhoneNumber: tenant.PhoneNumber,
			OpenHours:   tenant.OpenHours,
			LogoURL:     tenant.LogoURL,
			BannerURL:   tenant.BannerURL,
			IsOwner:     tenant.IsOwner,
		})
	}

	return &dto.ListPartnerTenantsResponse{Tenants: items}, nil
}

func (u *partnerTenantUsecase) SoftDeletePartnerTenant(ctx context.Context, userID string, tenantID string) error {
	if strings.TrimSpace(userID) == "" {
		return apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized)
	}
	if strings.TrimSpace(tenantID) == "" {
		return apperrors.New("VALIDATION_ERROR", "ID Tenant wajib diisi", http.StatusBadRequest)
	}

	err := u.repo.SoftDeleteTenant(ctx, userID, tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return apperrors.New("NOT_FOUND", "Tenant tidak ditemukan atau Anda tidak memiliki akses", http.StatusNotFound)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus tenant", http.StatusInternalServerError)
	}

	return nil
}

func (u *partnerTenantUsecase) GetPartnerTenantByID(ctx context.Context, userID string, tenantID string) (*dto.PartnerTenantResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized)
	}

	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "ID Tenant wajib diisi", http.StatusBadRequest)
	}

	tenant, err := u.repo.GetTenantByIDAndPartner(ctx, userID, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Tenant ID tidak ditemukan atau Anda tidak memiliki akses", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil detail tenant", http.StatusInternalServerError)
	}

	return &dto.PartnerTenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Slug:        tenant.Slug,
		Status:      tenant.Status,
		Description: tenant.Description,
		Address:     tenant.Address,
		PhoneNumber: tenant.PhoneNumber,
		OpenHours:   tenant.OpenHours,
		LogoURL:     tenant.LogoURL,
		BannerURL:   tenant.BannerURL,
		IsOwner:     tenant.IsOwner,
	}, nil
}

func (u *partnerTenantUsecase) UpdatePartnerTenant(ctx context.Context, userID string, tenantID string, req dto.UpdatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized)
	}

	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "ID Tenant wajib diisi", http.StatusBadRequest)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Nama tenant wajib diisi", http.StatusBadRequest)
	}

	tenant, err := u.repo.UpdateTenant(ctx, userID, tenantID, name, strings.TrimSpace(req.Description), strings.TrimSpace(req.Address), strings.TrimSpace(req.PhoneNumber))
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Tenant tidak ditemukan atau Anda tidak memiliki akses", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui tenant", http.StatusInternalServerError)
	}

	return &dto.CreatePartnerTenantResponse{
		Tenant: dto.PartnerTenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			Status:      tenant.Status,
			Description: tenant.Description,
			Address:     tenant.Address,
			PhoneNumber: tenant.PhoneNumber,
			OpenHours:   tenant.OpenHours,
			LogoURL:     tenant.LogoURL,
			BannerURL:   tenant.BannerURL,
			IsOwner:     tenant.IsOwner,
		},
	}, nil
}

