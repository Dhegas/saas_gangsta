package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/payment/dto"
)

// =====================================================
// Entities
// =====================================================

// OrderPaymentEntity representasi minimal order untuk proses webhook
type OrderPaymentEntity struct {
	ID                    string
	TenantID              string
	TotalPrice            float64
	PaymentStatus         string
	MidtransOrderID       *string
	MidtransTransactionID *string
	PaidAt                *time.Time
}

func (OrderPaymentEntity) TableName() string { return "orders" }

// TenantOwnerEntity untuk lookup user_id pemilik tenant
type TenantOwnerEntity struct {
	ID     string
	UserID string
}

func (TenantOwnerEntity) TableName() string { return "tenants" }

// PartnerWalletPaymentEntity untuk operasi wallet saat webhook
type PartnerWalletPaymentEntity struct {
	ID             string
	UserID         string
	Balance        float64
	TotalEarned    float64
	TotalWithdrawn float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (PartnerWalletPaymentEntity) TableName() string { return "partner_wallets" }

// WalletTransactionEntity untuk insert record transaksi
type WalletTransactionEntity struct {
	ID          string    `gorm:"primaryKey;default:gen_random_uuid()"`
	WalletID    string    `gorm:"not null"`
	OrderID     *string   `gorm:"index"`
	Type        string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	FeeAmount   float64   `gorm:"not null;default:0"`
	NetAmount   float64   `gorm:"not null"`
	Description string    `gorm:"not null;default:''"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (WalletTransactionEntity) TableName() string { return "wallet_transactions" }

// PlatformFeeConfigEntity untuk baca konfigurasi fee
type PlatformFeeConfigEntity struct {
	ID          string
	FeeType     string
	Amount      float64
	IsEnabled   bool
	Description string
}

func (PlatformFeeConfigEntity) TableName() string { return "platform_fee_configs" }

// =====================================================
// Interfaces
// =====================================================

// PaymentWebhookUsecase kontrak logika bisnis webhook
type PaymentWebhookUsecase interface {
	HandleMidtransWebhook(ctx context.Context, payload dto.MidtransWebhookPayload) error
	SyncPaymentStatus(ctx context.Context, orderID string) error
}

// PaymentSnapUsecase kontrak logika bisnis inisiasi pembayaran
type PaymentSnapUsecase interface {
	CreateSnapTransaction(ctx context.Context, userID, orderID string) (*dto.SnapResponse, error)
}

// PaymentRepository kontrak akses database untuk proses pembayaran
type PaymentRepository interface {
	// FindOrderByMidtransOrderID mencari order berdasarkan midtrans_order_id
	FindOrderByMidtransOrderID(ctx context.Context, midtransOrderID string) (*OrderPaymentEntity, error)

	// FindOrderByID mencari order berdasarkan ID untuk inisiasi Snap
	FindOrderByID(ctx context.Context, orderID string) (*OrderPaymentEntity, error)

	// GetOrderWithLockByID mencari order dengan lock FOR UPDATE
	GetOrderWithLockByID(ctx context.Context, orderID string) (*OrderPaymentEntity, error)

	// GetOrderWithLockByMidtransOrderID mencari order dengan lock FOR UPDATE
	GetOrderWithLockByMidtransOrderID(ctx context.Context, midtransOrderID string) (*OrderPaymentEntity, error)

	// GetTenantOwnerUserID mengambil user_id pemilik tenant
	GetTenantOwnerUserID(ctx context.Context, tenantID string) (string, error)

	// GetFeeConfig mengambil konfigurasi fee berdasarkan tipe
	GetFeeConfig(ctx context.Context, feeType string) (*PlatformFeeConfigEntity, error)

	// ProcessWebhookPayment menjalankan seluruh proses payment dalam satu DB transaction
	ProcessWebhookPayment(ctx context.Context, req ProcessPaymentRequest) error

	// UpdateMidtransOrderID menyimpan midtrans_order_id ke order saat inisiasi Snap
	UpdateMidtransOrderID(ctx context.Context, orderID, midtransOrderID string) error
}

// ProcessPaymentRequest data yang diperlukan untuk memproses webhook payment
type ProcessPaymentRequest struct {
	OrderID               string
	TenantID              string
	MidtransTransactionID string
	GrossAmount           float64
	FeeAmount             float64
	NetAmount             float64
	PaidAt                time.Time
}
