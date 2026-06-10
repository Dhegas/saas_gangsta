package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	walletdomain "github.com/dhegas/saas_gangsta/internal/domains/wallet/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/wallet/dto"
)

type adminWalletUsecase struct {
	repo walletdomain.WalletRepository
}

func NewAdminWalletUsecase(repo walletdomain.WalletRepository) walletdomain.AdminWalletUsecase {
	return &adminWalletUsecase{repo: repo}
}

// GetAllWithdrawals mengambil semua permintaan withdraw untuk admin
func (u *adminWalletUsecase) GetAllWithdrawals(ctx context.Context, filter walletdomain.WithdrawFilter) ([]dto.AdminWithdrawResponse, error) {
	withdrawals, err := u.repo.GetAllWithdrawals(ctx, filter)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil daftar penarikan", http.StatusInternalServerError)
	}

	result := make([]dto.AdminWithdrawResponse, 0, len(withdrawals))
	for _, w := range withdrawals {
		result = append(result, dto.AdminWithdrawResponse{
			ID:            w.ID,
			UserID:        w.UserID,
			WalletID:      w.WalletID,
			Amount:        w.Amount,
			FeeAmount:     w.FeeAmount,
			NetAmount:     w.NetAmount,
			Status:        w.Status,
			BankName:      w.BankName,
			BankAccount:   w.BankAccount,
			AccountHolder: w.AccountHolder,
			AdminNote:     w.AdminNote,
			ReviewedBy:    w.ReviewedBy,
			ReviewedAt:    w.ReviewedAt,
			TransferredAt: w.TransferredAt,
			CreatedAt:     w.CreatedAt,
			UpdatedAt:     w.UpdatedAt,
		})
	}
	return result, nil
}

// ApproveWithdrawal mengubah status withdraw menjadi APPROVED
func (u *adminWalletUsecase) ApproveWithdrawal(ctx context.Context, adminID, withdrawID string) error {
	if err := u.repo.ApproveWithdraw(ctx, withdrawID, adminID); err != nil {
		if strings.Contains(err.Error(), "INVALID_STATUS") {
			return apperrors.New("INVALID_STATUS", "Withdraw tidak ditemukan atau bukan status PENDING", http.StatusUnprocessableEntity)
		}
		return apperrors.New("APPROVE_FAILED", "Gagal menyetujui penarikan", http.StatusInternalServerError)
	}
	return nil
}

// RejectWithdrawal menolak withdraw dan mengembalikan saldo
func (u *adminWalletUsecase) RejectWithdrawal(ctx context.Context, adminID, withdrawID, note string) error {
	if err := u.repo.RejectWithdraw(ctx, withdrawID, adminID, note); err != nil {
		if strings.Contains(err.Error(), "INVALID_STATUS") {
			return apperrors.New("INVALID_STATUS", "Withdraw tidak ditemukan atau bukan status PENDING", http.StatusUnprocessableEntity)
		}
		return apperrors.New("REJECT_FAILED", "Gagal menolak penarikan", http.StatusInternalServerError)
	}
	return nil
}

// MarkTransferred menandai withdraw sebagai sudah ditransfer
func (u *adminWalletUsecase) MarkTransferred(ctx context.Context, adminID, withdrawID string) error {
	if err := u.repo.MarkAsTransferred(ctx, withdrawID, adminID); err != nil {
		if strings.Contains(err.Error(), "INVALID_STATUS") {
			return apperrors.New("INVALID_STATUS", "Withdraw tidak ditemukan atau sudah final", http.StatusUnprocessableEntity)
		}
		return apperrors.New("TRANSFER_FAILED", "Gagal menandai transfer", http.StatusInternalServerError)
	}
	return nil
}

// GetAllPartnerWallets mengambil semua wallet partner untuk admin
func (u *adminWalletUsecase) GetAllPartnerWallets(ctx context.Context, page, limit int) ([]dto.AdminPartnerWalletResponse, error) {
	wallets, err := u.repo.GetAllWallets(ctx, page, limit)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil daftar wallet partner", http.StatusInternalServerError)
	}

	result := make([]dto.AdminPartnerWalletResponse, 0, len(wallets))
	for _, w := range wallets {
		result = append(result, dto.AdminPartnerWalletResponse{
			WalletID:       w.ID,
			UserID:         w.UserID,
			Balance:        w.Balance,
			TotalEarned:    w.TotalEarned,
			TotalWithdrawn: w.TotalWithdrawn,
			UpdatedAt:      w.UpdatedAt,
		})
	}
	return result, nil
}
