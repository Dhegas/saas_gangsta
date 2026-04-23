package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/user/management/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/user/management/repository"
)

type UserUsecase interface {
	ListUsersByTenant(ctx context.Context, tenantID string) (*dto.ListUsersResponse, error)
	GetUserDetailByTenant(ctx context.Context, tenantID, userID string) (*dto.DetailUserResponse, error)
	UpdateUserByTenant(ctx context.Context, tenantID, userID string, req dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	SoftDeleteUserByTenant(ctx context.Context, tenantID, userID string) (*dto.DeleteUserResponse, error)
	ToggleUserActiveByTenant(ctx context.Context, tenantID, userID string) (*dto.ToggleActiveUserResponse, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) ListUsersByTenant(ctx context.Context, tenantID string) (*dto.ListUsersResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	users, err := u.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar user", http.StatusInternalServerError, nil)
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, dto.UserResponse{
			ID:       user.ID,
			TenantID: user.TenantID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			IsActive: user.IsActive,
		})
	}

	return &dto.ListUsersResponse{Users: items}, nil
}

func (u *userUsecase) GetUserDetailByTenant(ctx context.Context, tenantID, userID string) (*dto.DetailUserResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	user, err := u.repo.FindByIDAndTenant(ctx, tenantID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil detail user", http.StatusInternalServerError, nil)
	}

	return &dto.DetailUserResponse{User: dto.UserResponse{
		ID:       user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}}, nil
}

func (u *userUsecase) UpdateUserByTenant(ctx context.Context, tenantID, userID string, req dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	user, err := u.repo.UpdateByIDAndTenant(ctx, tenantID, userID, repository.UpdateUserInput{
		Email:    req.Email,
		FullName: req.FullName,
		Role:     req.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoFieldsToUpdate):
			return nil, apperrors.New("VALIDATION_ERROR", "Minimal satu field harus diisi untuk update user", http.StatusBadRequest, nil)
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound, nil)
		case errors.Is(err, repository.ErrEmailAlreadyExist):
			return nil, apperrors.New("CONFLICT", "Email sudah digunakan user lain", http.StatusConflict, nil)
		default:
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengupdate user", http.StatusInternalServerError, nil)
		}
	}

	return &dto.UpdateUserResponse{User: dto.UserResponse{
		ID:       user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}}, nil
}

func (u *userUsecase) SoftDeleteUserByTenant(ctx context.Context, tenantID, userID string) (*dto.DeleteUserResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	if err := u.repo.SoftDeleteByIDAndTenant(ctx, tenantID, userID); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menghapus user", http.StatusInternalServerError, nil)
	}

	return &dto.DeleteUserResponse{Deleted: true}, nil
}

func (u *userUsecase) ToggleUserActiveByTenant(ctx context.Context, tenantID, userID string) (*dto.ToggleActiveUserResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	user, err := u.repo.ToggleActiveByIDAndTenant(ctx, tenantID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengubah status aktif user", http.StatusInternalServerError, nil)
	}

	return &dto.ToggleActiveUserResponse{User: dto.UserResponse{
		ID:       user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		IsActive: user.IsActive,
	}}, nil
}
