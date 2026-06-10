package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/wallet/dto"
)

// =====================================================
// Entities
// =====================================================

// PartnerWalletEntity merepresentasikan tabel partner_wallets
type PartnerWalletEntity struct {
	ID             string    `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID         string    `gorm:"uniqueIndex;not null"`
	Balance        float64   `gorm:"not null;default:0"`
	TotalEarned    float64   `gorm:"not null;default:0"`
	TotalWithdrawn float64   `gorm:"not null;default:0"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (PartnerWalletEntity) TableName() string { return "partner_wallets" }

// WalletTransactionEntity merepresentasikan tabel wallet_transactions
type WalletTransactionEntity struct {
	ID          string    `gorm:"primaryKey;default:gen_random_uuid()"`
	WalletID    string    `gorm:"index;not null"`
	OrderID     *string   `gorm:"index"`
	Type        string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	FeeAmount   float64   `gorm:"not null;default:0"`
	NetAmount   float64   `gorm:"not null"`
	Description string    `gorm:"not null;default:''"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (WalletTransactionEntity) TableName() string { return "wallet_transactions" }

// WithdrawRequestEntity merepresentasikan tabel withdraw_requests
type WithdrawRequestEntity struct {
	ID            string     `gorm:"primaryKey;default:gen_random_uuid()"`
	WalletID      string     `gorm:"index;not null"`
	UserID        string     `gorm:"index;not null"`
	Amount        float64    `gorm:"not null"`
	FeeAmount     float64    `gorm:"not null;default:0"`
	NetAmount     float64    `gorm:"not null"`
	Status        string     `gorm:"not null;default:'PENDING'"`
	BankName      string     `gorm:"not null"`
	BankAccount   string     `gorm:"not null"`
	AccountHolder string     `gorm:"not null"`
	AdminNote     *string    ``
	ReviewedBy    *string    ``
	ReviewedAt    *time.Time ``
	TransferredAt *time.Time ``
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

func (WithdrawRequestEntity) TableName() string { return "withdraw_requests" }

// =====================================================
// Filter / Query types
// =====================================================

type WithdrawFilter struct {
	Status string
	Page   int
	Limit  int
}

type TransactionFilter struct {
	Page  int
	Limit int
}

// =====================================================
// Interfaces
// =====================================================

// PartnerWalletUsecase kontrak logika bisnis untuk partner
type PartnerWalletUsecase interface {
	GetWalletDashboard(ctx context.Context, userID string) (*dto.WalletDashboardResponse, error)
	GetTransactions(ctx context.Context, userID string, filter TransactionFilter) ([]dto.WalletTransactionResponse, error)
	CreateWithdrawRequest(ctx context.Context, userID string, req dto.CreateWithdrawRequest) (*dto.WithdrawResponse, error)
	GetMyWithdrawals(ctx context.Context, userID string, filter WithdrawFilter) ([]dto.WithdrawResponse, error)
	GetWithdrawalByID(ctx context.Context, userID, withdrawID string) (*dto.WithdrawResponse, error)
}

// AdminWalletUsecase kontrak logika bisnis untuk admin
type AdminWalletUsecase interface {
	GetAllWithdrawals(ctx context.Context, filter WithdrawFilter) ([]dto.AdminWithdrawResponse, error)
	ApproveWithdrawal(ctx context.Context, adminID, withdrawID string) error
	RejectWithdrawal(ctx context.Context, adminID, withdrawID, note string) error
	MarkTransferred(ctx context.Context, adminID, withdrawID string) error
	GetAllPartnerWallets(ctx context.Context, page, limit int) ([]dto.AdminPartnerWalletResponse, error)
}

// WalletRepository kontrak akses database untuk wallet
type WalletRepository interface {
	// Partner operations
	GetOrCreateWalletByUserID(ctx context.Context, userID string) (*PartnerWalletEntity, error)
	GetWalletByUserID(ctx context.Context, userID string) (*PartnerWalletEntity, error)
	GetTransactionsByWalletID(ctx context.Context, walletID string, filter TransactionFilter) ([]WalletTransactionEntity, error)
	GetPendingWithdrawSum(ctx context.Context, walletID string) (float64, error)

	// Withdraw operations (atomic)
	CreateWithdraw(ctx context.Context, walletID, userID string, req dto.CreateWithdrawRequest, feeAmount, netAmount float64) (*WithdrawRequestEntity, error)
	GetWithdrawsByUserID(ctx context.Context, userID string, filter WithdrawFilter) ([]WithdrawRequestEntity, error)
	GetWithdrawByID(ctx context.Context, withdrawID string) (*WithdrawRequestEntity, error)

	// Admin operations
	GetAllWithdrawals(ctx context.Context, filter WithdrawFilter) ([]WithdrawRequestEntity, error)
	ApproveWithdraw(ctx context.Context, withdrawID, adminID string) error
	RejectWithdraw(ctx context.Context, withdrawID, adminID, note string) error
	MarkAsTransferred(ctx context.Context, withdrawID, adminID string) error
	GetAllWallets(ctx context.Context, page, limit int) ([]PartnerWalletEntity, error)

	// Fee config
	GetWithdrawFeeConfig(ctx context.Context) (float64, bool, error)
}
