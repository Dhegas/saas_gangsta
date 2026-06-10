package usecase

import (
	"context"
	"crypto/sha512"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/config"
	paymentdomain "github.com/dhegas/saas_gangsta/internal/domains/payment/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/payment/dto"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentWebhookUsecase struct {
	repo paymentdomain.PaymentRepository
	cfg  *config.Config
}

type paymentSnapUsecase struct {
	repo paymentdomain.PaymentRepository
	cfg  *config.Config
}

func NewPaymentWebhookUsecase(repo paymentdomain.PaymentRepository, cfg *config.Config) paymentdomain.PaymentWebhookUsecase {
	return &paymentWebhookUsecase{repo: repo, cfg: cfg}
}

func NewPaymentSnapUsecase(repo paymentdomain.PaymentRepository, cfg *config.Config) paymentdomain.PaymentSnapUsecase {
	return &paymentSnapUsecase{repo: repo, cfg: cfg}
}

// =====================================================
// Webhook Usecase
// =====================================================

// HandleMidtransWebhook memproses notifikasi dari Midtrans secara idempotent.
// Flow:
//  1. Validasi signature (HMAC-SHA512)
//  2. Hanya proses status: settlement atau capture dengan fraud_status: accept
//  3. Lookup order by midtrans_order_id
//  4. Idempotency check: jika sudah PAID → return nil (sukses, jangan proses ulang)
//  5. Hitung fee transaksi (jika diaktifkan)
//  6. Proses pembayaran dalam DB transaction (atomic)
func (u *paymentWebhookUsecase) HandleMidtransWebhook(ctx context.Context, payload dto.MidtransWebhookPayload) error {
	// Step 1: Validasi signature Midtrans
	// Format: SHA512(order_id + status_code + gross_amount + ServerKey)
	if !u.validateSignature(payload) {
		slog.Warn("midtrans webhook signature invalid",
			slog.String("order_id", payload.OrderID),
			slog.String("status_code", payload.StatusCode),
		)
		return apperrors.New("INVALID_SIGNATURE", "Signature tidak valid", http.StatusForbidden)
	}

	// Step 2: Hanya proses transaksi yang berhasil
	// settlement = transfer bank/virtual account selesai
	// capture = kartu kredit/debit berhasil di-capture
	isSettled := payload.TransactionStatus == "settlement" || payload.TransactionStatus == "capture"
	isFraudClean := payload.FraudStatus == "accept" || payload.FraudStatus == ""
	if !isSettled || !isFraudClean {
		slog.Info("midtrans webhook skipped (not settlement or fraud challenge)",
			slog.String("order_id", payload.OrderID),
			slog.String("transaction_status", payload.TransactionStatus),
			slog.String("fraud_status", payload.FraudStatus),
		)
		return nil // Bukan error, hanya skip
	}

	// Step 3: Lookup order dengan Row-Level Locking (FOR UPDATE)
	order, err := u.repo.GetOrderWithLockByMidtransOrderID(ctx, payload.OrderID)
	if err != nil {
		slog.Error("midtrans webhook: order not found",
			slog.String("midtrans_order_id", payload.OrderID),
			slog.String("error", err.Error()),
		)
		return apperrors.New("ORDER_NOT_FOUND", "Order tidak ditemukan", http.StatusNotFound)
	}

	// Step 4: Idempotency Layer 1 — cek apakah order sudah PAID
	if order.PaymentStatus == "PAID" {
		slog.Info("midtrans webhook: order already paid, skipping",
			slog.String("order_id", order.ID),
		)
		return nil
	}

	// Step 5: Parse gross_amount dari string Midtrans
	grossAmount, err := strconv.ParseFloat(payload.GrossAmount, 64)
	if err != nil {
		return apperrors.New("INVALID_AMOUNT", "Gross amount tidak valid", http.StatusBadRequest)
	}

	// Step 6: Hitung fee transaksi (cek dari DB)
	feeAmount := 0.0
	feeConfig, err := u.repo.GetFeeConfig(ctx, "TRANSACTION")
	if err == nil && feeConfig.IsEnabled {
		feeAmount = feeConfig.Amount
	}
	netAmount := grossAmount - feeAmount
	if netAmount < 0 {
		netAmount = 0
	}

	// Step 7: Proses dalam DB transaction (atomic)
	processReq := paymentdomain.ProcessPaymentRequest{
		OrderID:               order.ID,
		TenantID:              order.TenantID,
		MidtransTransactionID: payload.TransactionID,
		GrossAmount:           grossAmount,
		FeeAmount:             feeAmount,
		NetAmount:             netAmount,
		PaidAt:                time.Now(),
	}

	if err := u.repo.ProcessWebhookPayment(ctx, processReq); err != nil {
		slog.Error("midtrans webhook: failed to process payment",
			slog.String("order_id", order.ID),
			slog.String("error", err.Error()),
		)
		return apperrors.New("PROCESS_FAILED", "Gagal memproses pembayaran", http.StatusInternalServerError)
	}

	slog.Info("midtrans webhook: payment processed successfully",
		slog.String("order_id", order.ID),
		slog.Float64("gross_amount", grossAmount),
		slog.Float64("fee_amount", feeAmount),
		slog.Float64("net_amount", netAmount),
	)
	return nil
}

// SyncPaymentStatus menyinkronkan status order lokal dengan status transaksi aktual di Midtrans secara langsung.
// Jika di Midtrans sukses (settlement/capture) tetapi lokal masih UNPAID, jalankan proses pembayaran.
func (u *paymentWebhookUsecase) SyncPaymentStatus(ctx context.Context, orderID string) error {
	// 1. Cari order di local database dengan Row-Level Locking (FOR UPDATE)
	order, err := u.repo.GetOrderWithLockByID(ctx, orderID)
	if err != nil {
		slog.Error("sync payment: order not found", slog.String("order_id", orderID), slog.String("error", err.Error()))
		return apperrors.New("ORDER_NOT_FOUND", "Order tidak ditemukan", http.StatusNotFound)
	}

	// 2. Cek idempotency: jika sudah PAID, skip
	if order.PaymentStatus == "PAID" {
		slog.Info("sync payment: order already paid locally", slog.String("order_id", order.ID))
		return nil
	}

	// 3. Tentukan midtrans_order_id
	midtransOrderID := "ORDER-" + order.ID
	if order.MidtransOrderID != nil && *order.MidtransOrderID != "" {
		midtransOrderID = *order.MidtransOrderID
	}

	// 4. Setup Midtrans CoreAPI client
	coreEnv := midtrans.Sandbox
	if u.cfg.MidtransEnv == "production" {
		coreEnv = midtrans.Production
	}

	var coreClient coreapi.Client
	coreClient.New(u.cfg.MidtransServerKey, coreEnv)

	// 5. Cek status transaksi langsung ke Midtrans API
	resp, snapErr := coreClient.CheckTransaction(midtransOrderID)
	if snapErr != nil {
		slog.Warn("sync payment: failed to check transaction from Midtrans",
			slog.String("midtrans_order_id", midtransOrderID),
			slog.String("error", snapErr.Error()),
		)
		return apperrors.New("MIDTRANS_API_ERROR", fmt.Sprintf("Gagal cek transaksi Midtrans: %s", snapErr.Message), http.StatusBadGateway)
	}

	// 6. Evaluasi status pembayaran
	isSettled := resp.TransactionStatus == "settlement" || resp.TransactionStatus == "capture"
	isFraudClean := resp.FraudStatus == "accept" || resp.FraudStatus == ""
	if !isSettled || !isFraudClean {
		slog.Info("sync payment: transaction is not settled yet in Midtrans",
			slog.String("order_id", order.ID),
			slog.String("transaction_status", resp.TransactionStatus),
			slog.String("fraud_status", resp.FraudStatus),
		)
		return apperrors.New("TRANSACTION_NOT_SETTLED", fmt.Sprintf("Transaksi di Midtrans belum selesai (status: %s)", resp.TransactionStatus), http.StatusUnprocessableEntity)
	}

	// 7. Parse gross_amount
	grossAmount, err := strconv.ParseFloat(resp.GrossAmount, 64)
	if err != nil {
		return apperrors.New("INVALID_AMOUNT", "Gross amount tidak valid", http.StatusBadRequest)
	}

	// 8. Hitung fee transaksi (sama dengan alur webhook)
	feeAmount := 0.0
	feeConfig, err := u.repo.GetFeeConfig(ctx, "TRANSACTION")
	if err == nil && feeConfig.IsEnabled {
		feeAmount = feeConfig.Amount
	}
	netAmount := grossAmount - feeAmount
	if netAmount < 0 {
		netAmount = 0
	}

	// 9. Jalankan proses pembayaran atomic (sama dengan alur webhook)
	processReq := paymentdomain.ProcessPaymentRequest{
		OrderID:               order.ID,
		TenantID:              order.TenantID,
		MidtransTransactionID: resp.TransactionID,
		GrossAmount:           grossAmount,
		FeeAmount:             feeAmount,
		NetAmount:             netAmount,
		PaidAt:                time.Now(),
	}

	if err := u.repo.ProcessWebhookPayment(ctx, processReq); err != nil {
		slog.Error("sync payment: failed to process payment recovery",
			slog.String("order_id", order.ID),
			slog.String("error", err.Error()),
		)
		return apperrors.New("PROCESS_FAILED", "Gagal menyinkronkan data pembayaran", http.StatusInternalServerError)
	}

	slog.Info("sync payment: status synchronized and updated successfully",
		slog.String("order_id", order.ID),
		slog.Float64("gross_amount", grossAmount),
		slog.Float64("fee_amount", feeAmount),
		slog.Float64("net_amount", netAmount),
	)
	return nil
}


// validateSignature memvalidasi signature dari Midtrans
// Formula: SHA512(order_id + status_code + gross_amount + ServerKey)
func (u *paymentWebhookUsecase) validateSignature(payload dto.MidtransWebhookPayload) bool {
	raw := payload.OrderID + payload.StatusCode + payload.GrossAmount + u.cfg.MidtransServerKey
	hash := sha512.Sum512([]byte(raw))
	expected := fmt.Sprintf("%x", hash)
	return strings.EqualFold(expected, payload.SignatureKey)
}

// =====================================================
// Snap Usecase
// =====================================================

// CreateSnapTransaction membuat Snap token Midtrans untuk order yang UNPAID.
// Token ini digunakan oleh front-end untuk menampilkan halaman pembayaran Midtrans.
func (u *paymentSnapUsecase) CreateSnapTransaction(ctx context.Context, userID, orderID string) (*dto.SnapResponse, error) {
	// 1. Cari order
	order, err := u.repo.FindOrderByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.New("ORDER_NOT_FOUND", "Order tidak ditemukan", http.StatusNotFound)
	}

	// 2. Hanya order UNPAID yang boleh dibuatkan payment
	if order.PaymentStatus != "UNPAID" {
		return nil, apperrors.New("ORDER_ALREADY_PAID", "Order sudah dibayar atau tidak valid", http.StatusUnprocessableEntity)
	}

	// 3. Gunakan midtrans_order_id yang sudah ada, atau gunakan order ID
	midtransOrderID := "ORDER-" + order.ID
	if order.MidtransOrderID != nil && *order.MidtransOrderID != "" {
		midtransOrderID = *order.MidtransOrderID
	}

	// 4. Simpan midtrans_order_id ke order jika belum ada
	if order.MidtransOrderID == nil || *order.MidtransOrderID == "" {
		if err := u.repo.UpdateMidtransOrderID(ctx, order.ID, midtransOrderID); err != nil {
			return nil, apperrors.New("UPDATE_FAILED", "Gagal menyimpan midtrans order ID", http.StatusInternalServerError)
		}
	}

	// 5. Setup Midtrans client
	snapEnv := midtrans.Sandbox
	if u.cfg.MidtransEnv == "production" {
		snapEnv = midtrans.Production
	}

	var snapClient snap.Client
	snapClient.New(u.cfg.MidtransServerKey, snapEnv)

	// 6. Buat Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  midtransOrderID,
			GrossAmt: int64(order.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{},
	}

	// 7. Hit Midtrans Snap API
	snapResp, snapErr := snapClient.CreateTransaction(snapReq)
	if snapErr != nil {
		slog.Error("midtrans snap: create transaction failed",
			slog.String("order_id", order.ID),
			slog.String("error", snapErr.Error()),
		)
		return nil, apperrors.New("MIDTRANS_ERROR", "Gagal membuat transaksi Midtrans", http.StatusBadGateway)
	}

	return &dto.SnapResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
		OrderID:     midtransOrderID,
		GrossAmount: order.TotalPrice,
	}, nil
}
