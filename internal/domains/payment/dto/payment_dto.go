package dto

// MidtransWebhookPayload adalah struktur payload yang dikirim Midtrans ke webhook kita
// Dokumentasi: https://docs.midtrans.com/docs/post-transaction-notification
type MidtransWebhookPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

// CreateSnapRequest request untuk inisiasi pembayaran Snap Midtrans
type CreateSnapRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

// SnapResponse response yang dikembalikan ke client setelah inisiasi Snap
type SnapResponse struct {
	Token       string `json:"snap_token"`
	RedirectURL string `json:"redirect_url"`
	OrderID     string `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}
