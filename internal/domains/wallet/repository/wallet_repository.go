package repository

import (
	"context"
	"fmt"
	"time"

	walletdomain "github.com/dhegas/saas_gangsta/internal/domains/wallet/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/wallet/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) walletdomain.WalletRepository {
	return &walletRepository{db: db}
}

// GetOrCreateWalletByUserID mengambil wallet partner, atau membuat baru jika belum ada
func (r *walletRepository) GetOrCreateWalletByUserID(ctx context.Context, userID string) (*walletdomain.PartnerWalletEntity, error) {
	var wallet walletdomain.PartnerWalletEntity

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err == gorm.ErrRecordNotFound {
		wallet = walletdomain.PartnerWalletEntity{
			UserID:  userID,
			Balance: 0,
		}
		if createErr := r.db.WithContext(ctx).Create(&wallet).Error; createErr != nil {
			return nil, createErr
		}
		return &wallet, nil
	}
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetWalletByUserID mengambil wallet partner berdasarkan user_id
func (r *walletRepository) GetWalletByUserID(ctx context.Context, userID string) (*walletdomain.PartnerWalletEntity, error) {
	var wallet walletdomain.PartnerWalletEntity
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetTransactionsByWalletID mengambil riwayat transaksi wallet dengan pagination
func (r *walletRepository) GetTransactionsByWalletID(ctx context.Context, walletID string, filter walletdomain.TransactionFilter) ([]walletdomain.WalletTransactionEntity, error) {
	var transactions []walletdomain.WalletTransactionEntity

	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, err
}

// GetPendingWithdrawSum mengambil total saldo yang sedang dalam proses withdraw (PENDING + APPROVED)
func (r *walletRepository) GetPendingWithdrawSum(ctx context.Context, walletID string) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Table("withdraw_requests").
		Select("COALESCE(SUM(amount), 0)").
		Where("wallet_id = ? AND status IN ('PENDING', 'APPROVED')", walletID).
		Scan(&total).Error
	return total, err
}

// CreateWithdraw membuat permintaan withdraw secara atomic:
// 1. Validasi saldo cukup (dengan lock)
// 2. Insert withdraw_request
// 3. Kurangi wallet.balance
// 4. Insert wallet_transaction DEBIT
func (r *walletRepository) CreateWithdraw(ctx context.Context, walletID, userID string, req dto.CreateWithdrawRequest, feeAmount, netAmount float64) (*walletdomain.WithdrawRequestEntity, error) {
	var created walletdomain.WithdrawRequestEntity

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Lock wallet row dan validasi saldo
		var wallet walletdomain.PartnerWalletEntity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", walletID).
			First(&wallet).Error; err != nil {
			return err
		}

		if wallet.Balance < req.Amount {
			return fmt.Errorf("INSUFFICIENT_BALANCE: saldo tidak mencukupi (saldo: %.2f, diminta: %.2f)", wallet.Balance, req.Amount)
		}

		// Step 2: Insert withdraw_request
		withdraw := walletdomain.WithdrawRequestEntity{
			WalletID:      walletID,
			UserID:        userID,
			Amount:        req.Amount,
			FeeAmount:     feeAmount,
			NetAmount:     netAmount,
			Status:        "PENDING",
			BankName:      req.BankName,
			BankAccount:   req.BankAccount,
			AccountHolder: req.AccountHolder,
		}
		if err := tx.Create(&withdraw).Error; err != nil {
			return err
		}
		created = withdraw

		// Step 3: Kurangi wallet.balance
		if err := tx.Table("partner_wallets").
			Where("id = ?", walletID).
			Updates(map[string]interface{}{
				"balance":    gorm.Expr("balance - ?", req.Amount),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		// Step 4: Insert wallet_transaction DEBIT
		if err := tx.Table("wallet_transactions").Create(map[string]interface{}{
			"wallet_id":   walletID,
			"order_id":    nil,
			"type":        "DEBIT",
			"amount":      req.Amount,
			"fee_amount":  feeAmount,
			"net_amount":  netAmount,
			"description": fmt.Sprintf("Penarikan saldo ke rekening %s - %s", req.BankName, req.BankAccount),
			"created_at":  time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &created, nil
}

// GetWithdrawsByUserID mengambil daftar withdraw milik partner dengan filter & pagination
func (r *walletRepository) GetWithdrawsByUserID(ctx context.Context, userID string, filter walletdomain.WithdrawFilter) ([]walletdomain.WithdrawRequestEntity, error) {
	var withdrawals []walletdomain.WithdrawRequestEntity

	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&withdrawals).Error

	return withdrawals, err
}

// GetWithdrawByID mengambil satu withdraw berdasarkan ID
func (r *walletRepository) GetWithdrawByID(ctx context.Context, withdrawID string) (*walletdomain.WithdrawRequestEntity, error) {
	var withdraw walletdomain.WithdrawRequestEntity
	err := r.db.WithContext(ctx).
		Where("id = ?", withdrawID).
		First(&withdraw).Error
	if err != nil {
		return nil, err
	}
	return &withdraw, nil
}

// GetAllWithdrawals mengambil semua withdraw untuk admin dengan filter & pagination
func (r *walletRepository) GetAllWithdrawals(ctx context.Context, filter walletdomain.WithdrawFilter) ([]walletdomain.WithdrawRequestEntity, error) {
	var withdrawals []walletdomain.WithdrawRequestEntity

	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := r.db.WithContext(ctx)
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&withdrawals).Error

	return withdrawals, err
}

// ApproveWithdraw mengubah status withdraw menjadi APPROVED
func (r *walletRepository) ApproveWithdraw(ctx context.Context, withdrawID, adminID string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Table("withdraw_requests").
		Where("id = ? AND status = 'PENDING'", withdrawID).
		Updates(map[string]interface{}{
			"status":      "APPROVED",
			"reviewed_by": adminID,
			"reviewed_at": &now,
			"updated_at":  now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("INVALID_STATUS: withdraw tidak ditemukan atau status bukan PENDING")
	}
	return nil
}

// RejectWithdraw menolak withdraw dan mengembalikan saldo ke wallet secara atomic
func (r *walletRepository) RejectWithdraw(ctx context.Context, withdrawID, adminID, note string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Ambil detail withdraw (harus status PENDING)
		var withdraw walletdomain.WithdrawRequestEntity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status = 'PENDING'", withdrawID).
			First(&withdraw).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("INVALID_STATUS: withdraw tidak ditemukan atau status bukan PENDING")
			}
			return err
		}

		// Update status menjadi REJECTED
		if err := tx.Table("withdraw_requests").
			Where("id = ?", withdrawID).
			Updates(map[string]interface{}{
				"status":      "REJECTED",
				"admin_note":  note,
				"reviewed_by": adminID,
				"reviewed_at": &now,
				"updated_at":  now,
			}).Error; err != nil {
			return err
		}

		// Kembalikan saldo ke wallet (REFUND)
		if err := tx.Table("partner_wallets").
			Where("id = ?", withdraw.WalletID).
			Updates(map[string]interface{}{
				"balance":    gorm.Expr("balance + ?", withdraw.Amount),
				"updated_at": now,
			}).Error; err != nil {
			return err
		}

		// Insert wallet_transaction CREDIT sebagai refund
		if err := tx.Table("wallet_transactions").Create(map[string]interface{}{
			"wallet_id":   withdraw.WalletID,
			"order_id":    nil,
			"type":        "CREDIT",
			"amount":      withdraw.Amount,
			"fee_amount":  0,
			"net_amount":  withdraw.Amount,
			"description": fmt.Sprintf("Refund penarikan ditolak - %s", note),
			"created_at":  now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// MarkAsTransferred menandai withdraw sebagai TRANSFERRED dan update total_withdrawn
func (r *walletRepository) MarkAsTransferred(ctx context.Context, withdrawID, adminID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Ambil detail withdraw (harus PENDING atau APPROVED)
		var withdraw walletdomain.WithdrawRequestEntity
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status IN ('PENDING', 'APPROVED')", withdrawID).
			First(&withdraw).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("INVALID_STATUS: withdraw tidak ditemukan atau sudah TRANSFERRED/REJECTED")
			}
			return err
		}

		// Update status
		if err := tx.Table("withdraw_requests").
			Where("id = ?", withdrawID).
			Updates(map[string]interface{}{
				"status":         "TRANSFERRED",
				"reviewed_by":    adminID,
				"reviewed_at":    &now,
				"transferred_at": &now,
				"updated_at":     now,
			}).Error; err != nil {
			return err
		}

		// Update total_withdrawn di wallet
		if err := tx.Table("partner_wallets").
			Where("id = ?", withdraw.WalletID).
			Updates(map[string]interface{}{
				"total_withdrawn": gorm.Expr("total_withdrawn + ?", withdraw.Amount),
				"updated_at":      now,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAllWallets mengambil semua wallet partner untuk admin
func (r *walletRepository) GetAllWallets(ctx context.Context, page, limit int) ([]walletdomain.PartnerWalletEntity, error) {
	var wallets []walletdomain.PartnerWalletEntity

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	err := r.db.WithContext(ctx).
		Order("total_earned DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&wallets).Error

	return wallets, err
}

// GetWithdrawFeeConfig mengambil konfigurasi fee withdraw dari DB
func (r *walletRepository) GetWithdrawFeeConfig(ctx context.Context) (float64, bool, error) {
	var config struct {
		Amount    float64 `gorm:"column:amount"`
		IsEnabled bool    `gorm:"column:is_enabled"`
	}

	err := r.db.WithContext(ctx).
		Table("platform_fee_configs").
		Select("amount, is_enabled").
		Where("fee_type = 'WITHDRAW'").
		First(&config).Error
	if err != nil {
		// Default: fee Rp2000, enabled
		return 2000, true, err
	}

	return config.Amount, config.IsEnabled, nil
}
