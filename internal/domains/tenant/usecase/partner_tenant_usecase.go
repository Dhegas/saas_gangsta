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
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Nama tenant wajib diisi", http.StatusBadRequest, nil)
	}

	tenant, err := u.repo.CreateTenantForPartner(ctx, domain.CreatePartnerTenantInput{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound, nil)
		case errors.Is(err, repository.ErrUserNotPartner):
			return nil, apperrors.New("FORBIDDEN", "Hanya PARTNER yang dapat membuat tenant", http.StatusForbidden, nil)
		case errors.Is(err, repository.ErrPartnerSubscriptionMissing):
			return nil, apperrors.New("FORBIDDEN", "Subscription PARTNER tidak ditemukan", http.StatusForbidden, nil)
		case errors.Is(err, repository.ErrTenantLimitReached):
			return nil, apperrors.New("FORBIDDEN", "Batas jumlah tenant pada paket subscription sudah tercapai", http.StatusForbidden, nil)
		default:
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat tenant partner", http.StatusInternalServerError, nil)
		}
	}

	return &dto.CreatePartnerTenantResponse{
		Tenant: dto.PartnerTenantResponse{
			ID:      tenant.ID,
			Name:    tenant.Name,
			Slug:    tenant.Slug,
			Status:  tenant.Status,
			IsOwner: tenant.IsOwner,
		},
	}, nil
}

func (u *partnerTenantUsecase) ListPartnerTenants(ctx context.Context, userID string) (*dto.ListPartnerTenantsResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	partner, err := u.repo.FindPartnerByID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data user", http.StatusInternalServerError, nil)
	}
	if partner == nil || !partner.IsActive {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak ditemukan", http.StatusUnauthorized, nil)
	}
	if partner.Role != "PARTNER" {
		return nil, apperrors.New("FORBIDDEN", "Hanya PARTNER yang dapat melihat tenant", http.StatusForbidden, nil)
	}

	tenants, err := u.repo.ListTenantsByPartner(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar tenant partner", http.StatusInternalServerError, nil)
	}

	items := make([]dto.PartnerTenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		items = append(items, dto.PartnerTenantResponse{
			ID:      tenant.ID,
			Name:    tenant.Name,
			Slug:    tenant.Slug,
			Status:  tenant.Status,
			IsOwner: tenant.IsOwner,
		})
	}

	return &dto.ListPartnerTenantsResponse{Tenants: items}, nil
}
