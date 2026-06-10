package dto

import "time"

// =====================================================
// Partner DTOs
// =====================================================

// WalletDashboardResponse data wallet yang ditampilkan di dashboard partner
type WalletDashboardResponse struct {
	WalletID           string    `json:"wallet_id"`
	Balance            float64   `json:"balance"`              // Saldo tersedia untuk withdraw
	WithdrawInProgress float64   `json:"withdraw_in_progress"` // Total PENDING + APPROVED withdraw
	TotalEarned        float64   `json:"total_earned"`         // Total seluruh pendapatan
	TotalWithdrawn     float64   `json:"total_withdrawn"`      // Total withdraw yang sudah TRANSFERRED
	UpdatedAt          time.Time `json:"updated_at"`
}

// WalletTransactionResponse representasi satu entri riwayat transaksi wallet
type WalletTransactionResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`        // CREDIT atau DEBIT
	Amount      float64   `json:"amount"`      // Jumlah gross
	FeeAmount   float64   `json:"fee_amount"`  // Fee yang dipotong
	NetAmount   float64   `json:"net_amount"`  // Jumlah bersih
	Description string    `json:"description"`
	OrderID     *string   `json:"order_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateWithdrawRequest request untuk membuat permintaan penarikan saldo
type CreateWithdrawRequest struct {
	Amount        float64 `json:"amount"         binding:"required,gt=0"`
	BankName      string  `json:"bank_name"      binding:"required,min=2,max=100"`
	BankAccount   string  `json:"bank_account"   binding:"required,min=5,max=50"`
	AccountHolder string  `json:"account_holder" binding:"required,min=2,max=150"`
}

// WithdrawResponse response detail permintaan withdraw
type WithdrawResponse struct {
	ID            string     `json:"id"`
	Amount        float64    `json:"amount"`         // Jumlah yang diminta (sebelum fee)
	FeeAmount     float64    `json:"fee_amount"`     // Biaya admin withdraw
	NetAmount     float64    `json:"net_amount"`     // Jumlah yang ditransfer ke rekening
	Status        string     `json:"status"`
	BankName      string     `json:"bank_name"`
	BankAccount   string     `json:"bank_account"`
	AccountHolder string     `json:"account_holder"`
	AdminNote     *string    `json:"admin_note,omitempty"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	TransferredAt *time.Time `json:"transferred_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// =====================================================
// Admin DTOs
// =====================================================

// AdminWithdrawResponse response untuk admin melihat daftar withdraw
type AdminWithdrawResponse struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	WalletID      string     `json:"wallet_id"`
	Amount        float64    `json:"amount"`
	FeeAmount     float64    `json:"fee_amount"`
	NetAmount     float64    `json:"net_amount"`
	Status        string     `json:"status"`
	BankName      string     `json:"bank_name"`
	BankAccount   string     `json:"bank_account"`
	AccountHolder string     `json:"account_holder"`
	AdminNote     *string    `json:"admin_note,omitempty"`
	ReviewedBy    *string    `json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	TransferredAt *time.Time `json:"transferred_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AdminPartnerWalletResponse ringkasan wallet setiap partner untuk view admin
type AdminPartnerWalletResponse struct {
	WalletID       string    `json:"wallet_id"`
	UserID         string    `json:"user_id"`
	Balance        float64   `json:"balance"`
	TotalEarned    float64   `json:"total_earned"`
	TotalWithdrawn float64   `json:"total_withdrawn"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AdminReviewRequest body untuk aksi reject dari admin (berisi admin_note)
type AdminReviewRequest struct {
	AdminNote string `json:"admin_note" binding:"required,min=3"`
}
