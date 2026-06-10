package usecase

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	walletdomain "github.com/dhegas/saas_gangsta/internal/domains/wallet/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/wallet/dto"
)

const (
	MinWithdrawAmount = 20000.0 // Minimum withdraw Rp20.000
)

type partnerWalletUsecase struct {
	repo walletdomain.WalletRepository
}

func NewPartnerWalletUsecase(repo walletdomain.WalletRepository) walletdomain.PartnerWalletUsecase {
	return &partnerWalletUsecase{repo: repo}
}

// GetWalletDashboard mengambil ringkasan wallet partner untuk dashboard
func (u *partnerWalletUsecase) GetWalletDashboard(ctx context.Context, userID string) (*dto.WalletDashboardResponse, error) {
	wallet, err := u.repo.GetOrCreateWalletByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil data wallet", http.StatusInternalServerError)
	}

	// Hitung total yang sedang dalam proses withdraw
	pendingSum, err := u.repo.GetPendingWithdrawSum(ctx, wallet.ID)
	if err != nil {
		pendingSum = 0 // Non-fatal
	}

	return &dto.WalletDashboardResponse{
		WalletID:           wallet.ID,
		Balance:            wallet.Balance,
		WithdrawInProgress: pendingSum,
		TotalEarned:        wallet.TotalEarned,
		TotalWithdrawn:     wallet.TotalWithdrawn,
		UpdatedAt:          wallet.UpdatedAt,
	}, nil
}

// GetTransactions mengambil riwayat transaksi wallet partner
func (u *partnerWalletUsecase) GetTransactions(ctx context.Context, userID string, filter walletdomain.TransactionFilter) ([]dto.WalletTransactionResponse, error) {
	wallet, err := u.repo.GetOrCreateWalletByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil data wallet", http.StatusInternalServerError)
	}

	transactions, err := u.repo.GetTransactionsByWalletID(ctx, wallet.ID, filter)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil riwayat transaksi", http.StatusInternalServerError)
	}

	result := make([]dto.WalletTransactionResponse, 0, len(transactions))
	for _, tx := range transactions {
		result = append(result, dto.WalletTransactionResponse{
			ID:          tx.ID,
			Type:        tx.Type,
			Amount:      tx.Amount,
			FeeAmount:   tx.FeeAmount,
			NetAmount:   tx.NetAmount,
			Description: tx.Description,
			OrderID:     tx.OrderID,
			CreatedAt:   tx.CreatedAt,
		})
	}
	return result, nil
}

// CreateWithdrawRequest membuat permintaan penarikan saldo.
// Validasi:
//   - amount >= MinWithdrawAmount (Rp20.000)
//   - amount <= wallet.balance
//   - Fee withdraw dikurangi dari net (bukan dari balance langsung)
func (u *partnerWalletUsecase) CreateWithdrawRequest(ctx context.Context, userID string, req dto.CreateWithdrawRequest) (*dto.WithdrawResponse, error) {
	// Validasi minimum withdraw
	if req.Amount < MinWithdrawAmount {
		return nil, apperrors.New(
			"AMOUNT_TOO_SMALL",
			fmt.Sprintf("Jumlah penarikan minimum adalah Rp%.0f", MinWithdrawAmount),
			http.StatusUnprocessableEntity,
		)
	}

	// Ambil wallet
	wallet, err := u.repo.GetOrCreateWalletByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil data wallet", http.StatusInternalServerError)
	}

	// Validasi saldo
	if wallet.Balance < req.Amount {
		return nil, apperrors.New(
			"INSUFFICIENT_BALANCE",
			fmt.Sprintf("Saldo tidak mencukupi. Saldo tersedia: Rp%.0f", wallet.Balance),
			http.StatusUnprocessableEntity,
		)
	}

	// Ambil konfigurasi fee withdraw
	feeAmount := 0.0
	feeValue, feeEnabled, _ := u.repo.GetWithdrawFeeConfig(ctx)
	if feeEnabled {
		feeAmount = feeValue
	}
	netAmount := req.Amount - feeAmount
	if netAmount < 0 {
		netAmount = 0
	}

	// Buat withdraw secara atomic
	withdraw, err := u.repo.CreateWithdraw(ctx, wallet.ID, userID, req, feeAmount, netAmount)
	if err != nil {
		if strings.Contains(err.Error(), "INSUFFICIENT_BALANCE") {
			return nil, apperrors.New("INSUFFICIENT_BALANCE", "Saldo tidak mencukupi", http.StatusUnprocessableEntity)
		}
		return nil, apperrors.New("WITHDRAW_FAILED", "Gagal membuat permintaan penarikan", http.StatusInternalServerError)
	}

	return mapWithdrawToResponse(withdraw), nil
}

// GetMyWithdrawals mengambil daftar withdraw milik partner
func (u *partnerWalletUsecase) GetMyWithdrawals(ctx context.Context, userID string, filter walletdomain.WithdrawFilter) ([]dto.WithdrawResponse, error) {
	withdrawals, err := u.repo.GetWithdrawsByUserID(ctx, userID, filter)
	if err != nil {
		return nil, apperrors.New("WALLET_ERROR", "Gagal mengambil riwayat penarikan", http.StatusInternalServerError)
	}

	result := make([]dto.WithdrawResponse, 0, len(withdrawals))
	for i := range withdrawals {
		result = append(result, *mapWithdrawToResponse(&withdrawals[i]))
	}
	return result, nil
}

// GetWithdrawalByID mengambil detail satu withdraw milik partner
// Validasi IDOR: withdraw.user_id harus sama dengan userID dari JWT
func (u *partnerWalletUsecase) GetWithdrawalByID(ctx context.Context, userID, withdrawID string) (*dto.WithdrawResponse, error) {
	withdraw, err := u.repo.GetWithdrawByID(ctx, withdrawID)
	if err != nil {
		return nil, apperrors.New("NOT_FOUND", "Permintaan penarikan tidak ditemukan", http.StatusNotFound)
	}

	// Proteksi IDOR
	if withdraw.UserID != userID {
		return nil, apperrors.New("FORBIDDEN", "Anda tidak memiliki akses ke permintaan ini", http.StatusForbidden)
	}

	return mapWithdrawToResponse(withdraw), nil
}

// mapWithdrawToResponse mengkonversi entity ke DTO
func mapWithdrawToResponse(w *walletdomain.WithdrawRequestEntity) *dto.WithdrawResponse {
	return &dto.WithdrawResponse{
		ID:            w.ID,
		Amount:        w.Amount,
		FeeAmount:     w.FeeAmount,
		NetAmount:     w.NetAmount,
		Status:        w.Status,
		BankName:      w.BankName,
		BankAccount:   w.BankAccount,
		AccountHolder: w.AccountHolder,
		AdminNote:     w.AdminNote,
		ReviewedAt:    w.ReviewedAt,
		TransferredAt: w.TransferredAt,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}
