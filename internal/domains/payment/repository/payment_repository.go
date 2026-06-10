package repository

import (
	"context"
	"time"

	paymentdomain "github.com/dhegas/saas_gangsta/internal/domains/payment/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) paymentdomain.PaymentRepository {
	return &paymentRepository{db: db}
}

// FindOrderByMidtransOrderID mencari order berdasarkan midtrans_order_id
func (r *paymentRepository) FindOrderByMidtransOrderID(ctx context.Context, midtransOrderID string) (*paymentdomain.OrderPaymentEntity, error) {
	var order struct {
		ID                    string     `gorm:"column:id"`
		TenantID              string     `gorm:"column:tenant_id"`
		TotalPrice            float64    `gorm:"column:total_price"`
		PaymentStatus         string     `gorm:"column:payment_status"`
		MidtransOrderID       *string    `gorm:"column:midtrans_order_id"`
		MidtransTransactionID *string    `gorm:"column:midtrans_transaction_id"`
		PaidAt                *time.Time `gorm:"column:paid_at"`
	}

	err := r.db.WithContext(ctx).
		Table("orders").
		Select("id, tenant_id, total_price, payment_status, midtrans_order_id, midtrans_transaction_id, paid_at").
		Where("midtrans_order_id = ? AND deleted_at IS NULL", midtransOrderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return &paymentdomain.OrderPaymentEntity{
		ID:                    order.ID,
		TenantID:              order.TenantID,
		TotalPrice:            order.TotalPrice,
		PaymentStatus:         order.PaymentStatus,
		MidtransOrderID:       order.MidtransOrderID,
		MidtransTransactionID: order.MidtransTransactionID,
		PaidAt:                order.PaidAt,
	}, nil
}

// FindOrderByID mencari order berdasarkan ID (untuk inisiasi Snap)
func (r *paymentRepository) FindOrderByID(ctx context.Context, orderID string) (*paymentdomain.OrderPaymentEntity, error) {
	var order struct {
		ID                    string     `gorm:"column:id"`
		TenantID              string     `gorm:"column:tenant_id"`
		TotalPrice            float64    `gorm:"column:total_price"`
		PaymentStatus         string     `gorm:"column:payment_status"`
		MidtransOrderID       *string    `gorm:"column:midtrans_order_id"`
		MidtransTransactionID *string    `gorm:"column:midtrans_transaction_id"`
		PaidAt                *time.Time `gorm:"column:paid_at"`
	}

	err := r.db.WithContext(ctx).
		Table("orders").
		Select("id, tenant_id, total_price, payment_status, midtrans_order_id, midtrans_transaction_id, paid_at").
		Where("id = ? AND deleted_at IS NULL", orderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return &paymentdomain.OrderPaymentEntity{
		ID:                    order.ID,
		TenantID:              order.TenantID,
		TotalPrice:            order.TotalPrice,
		PaymentStatus:         order.PaymentStatus,
		MidtransOrderID:       order.MidtransOrderID,
		MidtransTransactionID: order.MidtransTransactionID,
		PaidAt:                order.PaidAt,
	}, nil
}

// GetTenantOwnerUserID mengambil user_id pemilik tenant
func (r *paymentRepository) GetTenantOwnerUserID(ctx context.Context, tenantID string) (string, error) {
	var result struct {
		UserID string `gorm:"column:user_id"`
	}

	err := r.db.WithContext(ctx).
		Table("tenants").
		Select("user_id").
		Where("id = ? AND deleted_at IS NULL", tenantID).
		First(&result).Error
	if err != nil {
		return "", err
	}
	return result.UserID, nil
}

// GetFeeConfig mengambil konfigurasi fee berdasarkan tipe
func (r *paymentRepository) GetFeeConfig(ctx context.Context, feeType string) (*paymentdomain.PlatformFeeConfigEntity, error) {
	var config struct {
		ID          string  `gorm:"column:id"`
		FeeType     string  `gorm:"column:fee_type"`
		Amount      float64 `gorm:"column:amount"`
		IsEnabled   bool    `gorm:"column:is_enabled"`
		Description string  `gorm:"column:description"`
	}

	err := r.db.WithContext(ctx).
		Table("platform_fee_configs").
		Where("fee_type = ?", feeType).
		First(&config).Error
	if err != nil {
		return nil, err
	}

	return &paymentdomain.PlatformFeeConfigEntity{
		ID:          config.ID,
		FeeType:     config.FeeType,
		Amount:      config.Amount,
		IsEnabled:   config.IsEnabled,
		Description: config.Description,
	}, nil
}

// ProcessWebhookPayment menjalankan seluruh proses pembayaran dalam satu DB transaction.
// Ini adalah operasi utama dan harus ATOMIC:
//  1. Update order payment_status = PAID
//  2. Get or create partner wallet (FOR UPDATE untuk cegah race condition)
//  3. Credit wallet balance + total_earned
//  4. Insert wallet_transaction CREDIT (constraint unique order_id menjamin idempotency layer 2)
func (r *paymentRepository) ProcessWebhookPayment(ctx context.Context, req paymentdomain.ProcessPaymentRequest) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Step 1: Update order menjadi PAID
		txResult := tx.
			Table("orders").
			Where("id = ? AND payment_status = 'UNPAID'", req.OrderID).
			Updates(map[string]interface{}{
				"payment_status":           "PAID",
				"midtrans_transaction_id":  req.MidtransTransactionID,
				"paid_at":                  &now,
				"updated_at":               now,
			})
		if txResult.Error != nil {
			return txResult.Error
		}
		// Jika rows tidak terpengaruh, order sudah PAID (idempotent — sukses)
		if txResult.RowsAffected == 0 {
			return nil
		}

		// Step 2: Cari user_id pemilik tenant
		var tenantResult struct {
			UserID string `gorm:"column:user_id"`
		}
		if err := tx.Table("tenants").
			Select("user_id").
			Where("id = ?", req.TenantID).
			First(&tenantResult).Error; err != nil {
			return err
		}

		// Step 3: GET OR CREATE wallet partner + LOCK row untuk cegah race condition
		wallet := struct {
			ID          string  `gorm:"column:id"`
			Balance     float64 `gorm:"column:balance"`
			TotalEarned float64 `gorm:"column:total_earned"`
		}{}

		err := tx.Table("partner_wallets").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", tenantResult.UserID).
			First(&wallet).Error

		if err == gorm.ErrRecordNotFound {
			// Wallet belum ada — buat baru
			if createErr := tx.Table("partner_wallets").Create(map[string]interface{}{
				"user_id":         tenantResult.UserID,
				"balance":         req.NetAmount,
				"total_earned":    req.NetAmount,
				"total_withdrawn": 0.00,
				"created_at":      now,
				"updated_at":      now,
			}).Error; createErr != nil {
				return createErr
			}

			// Ambil wallet yang baru dibuat untuk dapatkan ID
			if err2 := tx.Table("partner_wallets").
				Where("user_id = ?", tenantResult.UserID).
				First(&wallet).Error; err2 != nil {
				return err2
			}
		} else if err != nil {
			return err
		} else {
			// Step 4: Update balance + total_earned yang sudah ada
			if err := tx.Table("partner_wallets").
				Where("id = ?", wallet.ID).
				Updates(map[string]interface{}{
					"balance":      gorm.Expr("balance + ?", req.NetAmount),
					"total_earned": gorm.Expr("total_earned + ?", req.NetAmount),
					"updated_at":   now,
				}).Error; err != nil {
				return err
			}
		}

		// Step 5: Insert wallet_transaction CREDIT
		// Jika order_id sudah ada (UNIQUE constraint), ini akan gagal → idempotency layer 2
		orderID := req.OrderID
		if err := tx.Table("wallet_transactions").Create(map[string]interface{}{
			"wallet_id":   wallet.ID,
			"order_id":    &orderID,
			"type":        "CREDIT",
			"amount":      req.GrossAmount,
			"fee_amount":  req.FeeAmount,
			"net_amount":  req.NetAmount,
			"description": "Pendapatan dari order #" + req.OrderID,
			"created_at":  now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateMidtransOrderID menyimpan midtrans_order_id ke order saat inisiasi Snap
func (r *paymentRepository) UpdateMidtransOrderID(ctx context.Context, orderID, midtransOrderID string) error {
	return r.db.WithContext(ctx).
		Table("orders").
		Where("id = ?", orderID).
		Update("midtrans_order_id", midtransOrderID).Error
}

// GetOrderWithLockByID mencari order dengan lock FOR UPDATE dalam transaksi singkat
func (r *paymentRepository) GetOrderWithLockByID(ctx context.Context, orderID string) (*paymentdomain.OrderPaymentEntity, error) {
	var order struct {
		ID                    string     `gorm:"column:id"`
		TenantID              string     `gorm:"column:tenant_id"`
		TotalPrice            float64    `gorm:"column:total_price"`
		PaymentStatus         string     `gorm:"column:payment_status"`
		MidtransOrderID       *string    `gorm:"column:midtrans_order_id"`
		MidtransTransactionID *string    `gorm:"column:midtrans_transaction_id"`
		PaidAt                *time.Time `gorm:"column:paid_at"`
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Table("orders").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, tenant_id, total_price, payment_status, midtrans_order_id, midtrans_transaction_id, paid_at").
			Where("id = ? AND deleted_at IS NULL", orderID).
			First(&order).Error
	})
	if err != nil {
		return nil, err
	}

	return &paymentdomain.OrderPaymentEntity{
		ID:                    order.ID,
		TenantID:              order.TenantID,
		TotalPrice:            order.TotalPrice,
		PaymentStatus:         order.PaymentStatus,
		MidtransOrderID:       order.MidtransOrderID,
		MidtransTransactionID: order.MidtransTransactionID,
		PaidAt:                order.PaidAt,
	}, nil
}

// GetOrderWithLockByMidtransOrderID mencari order dengan lock FOR UPDATE dalam transaksi singkat
func (r *paymentRepository) GetOrderWithLockByMidtransOrderID(ctx context.Context, midtransOrderID string) (*paymentdomain.OrderPaymentEntity, error) {
	var order struct {
		ID                    string     `gorm:"column:id"`
		TenantID              string     `gorm:"column:tenant_id"`
		TotalPrice            float64    `gorm:"column:total_price"`
		PaymentStatus         string     `gorm:"column:payment_status"`
		MidtransOrderID       *string    `gorm:"column:midtrans_order_id"`
		MidtransTransactionID *string    `gorm:"column:midtrans_transaction_id"`
		PaidAt                *time.Time `gorm:"column:paid_at"`
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Table("orders").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, tenant_id, total_price, payment_status, midtrans_order_id, midtrans_transaction_id, paid_at").
			Where("midtrans_order_id = ? AND deleted_at IS NULL", midtransOrderID).
			First(&order).Error
	})
	if err != nil {
		return nil, err
	}

	return &paymentdomain.OrderPaymentEntity{
		ID:                    order.ID,
		TenantID:              order.TenantID,
		TotalPrice:            order.TotalPrice,
		PaymentStatus:         order.PaymentStatus,
		MidtransOrderID:       order.MidtransOrderID,
		MidtransTransactionID: order.MidtransTransactionID,
		PaidAt:                order.PaidAt,
	}, nil
}
