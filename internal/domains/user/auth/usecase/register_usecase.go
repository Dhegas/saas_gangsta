package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)
	fullName := strings.TrimSpace(req.FullName)
	role := strings.ToUpper(strings.TrimSpace(req.Role))

	// Jika role tidak dikirimkan, gunakan default "CUSTOMER"
	if role == "" {
		role = "CUSTOMER"
	}

	// Validasi role yang diperbolehkan untuk registrasi publik.
	// Membatasi registrasi ADMIN secara publik untuk alasan keamanan.
	if role != "CUSTOMER" && role != "PARTNER" {
		if role == "ADMIN" {
			return nil, apperrors.New("VALIDATION_ERROR", "Registrasi sebagai ADMIN tidak diperbolehkan secara publik", http.StatusBadRequest)
		}
		return nil, apperrors.New("VALIDATION_ERROR", "Role tidak valid. Harus salah satu dari: CUSTOMER atau PARTNER", http.StatusBadRequest)
	}

	existing, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memproses registrasi", http.StatusInternalServerError)
	}
	if existing != nil {
		return nil, apperrors.New("CONFLICT", "Email sudah terdaftar", http.StatusConflict)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyiapkan akun", http.StatusInternalServerError)
	}

	user := &domain.User{
		TenantID:     "",
		Email:        email,
		FullName:     fullName,
		PasswordHash: string(passwordHash),
		Role:         role,
		IsActive:     true,
		TenantStatus: "active",
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate email") {
			return nil, apperrors.New("CONFLICT", "Email sudah terdaftar", http.StatusConflict)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan akun", http.StatusInternalServerError)
	}

	return &dto.RegisterResponse{
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}
