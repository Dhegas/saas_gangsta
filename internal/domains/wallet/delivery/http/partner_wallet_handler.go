package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	walletdomain "github.com/dhegas/saas_gangsta/internal/domains/wallet/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/wallet/dto"
)

// PartnerWalletHandler handler untuk endpoint wallet partner
type PartnerWalletHandler struct {
	usecase walletdomain.PartnerWalletUsecase
}

func NewPartnerWalletHandler(usecase walletdomain.PartnerWalletUsecase) *PartnerWalletHandler {
	return &PartnerWalletHandler{usecase: usecase}
}

// GetWalletDashboard godoc
// @Summary      Dashboard Wallet Partner
// @Description  Mengambil ringkasan saldo, total pendapatan, dan total withdraw partner yang sedang login
// @Tags         Partner Wallet
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /partner/wallet [get]
func (h *PartnerWalletHandler) GetWalletDashboard(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "User tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	result, err := h.usecase.GetWalletDashboard(c.Request.Context(), userID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Data wallet berhasil diambil", result)
}

// GetTransactions godoc
// @Summary      Riwayat Transaksi Wallet
// @Description  Mengambil riwayat CREDIT dan DEBIT wallet partner dengan pagination
// @Tags         Partner Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int  false  "Halaman (default: 1)"
// @Param        limit  query  int  false  "Jumlah per halaman (default: 20, max: 100)"
// @Success      200  {object}  map[string]interface{}
// @Router       /partner/wallet/transactions [get]
func (h *PartnerWalletHandler) GetTransactions(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "User tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := walletdomain.TransactionFilter{Page: page, Limit: limit}
	result, err := h.usecase.GetTransactions(c.Request.Context(), userID, filter)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Riwayat transaksi berhasil diambil", result)
}

// CreateWithdraw godoc
// @Summary      Buat Permintaan Penarikan Saldo
// @Description  Partner membuat permintaan withdraw. Saldo langsung dikurangi dan masuk status PENDING. Minimum withdraw Rp20.000. Fee withdraw Rp2.000 dipotong dari net amount.
// @Tags         Partner Wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  dto.CreateWithdrawRequest  true  "Data penarikan"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /partner/wallet/withdraw [post]
func (h *PartnerWalletHandler) CreateWithdraw(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "User tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	var req dto.CreateWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validasi gagal", gin.H{
			"code":  "VALIDATION_ERROR",
			"error": err.Error(),
		})
		return
	}

	result, err := h.usecase.CreateWithdrawRequest(c.Request.Context(), userID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Permintaan penarikan berhasil dibuat", result)
}

// GetMyWithdrawals godoc
// @Summary      Daftar Withdraw Saya
// @Description  Mengambil semua permintaan withdraw partner yang sedang login
// @Tags         Partner Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        status  query  string  false  "Filter status: PENDING, APPROVED, TRANSFERRED, REJECTED"
// @Param        page    query  int     false  "Halaman (default: 1)"
// @Param        limit   query  int     false  "Jumlah per halaman (default: 20)"
// @Success      200  {object}  map[string]interface{}
// @Router       /partner/wallet/withdrawals [get]
func (h *PartnerWalletHandler) GetMyWithdrawals(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "User tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	filter := walletdomain.WithdrawFilter{Status: status, Page: page, Limit: limit}
	result, err := h.usecase.GetMyWithdrawals(c.Request.Context(), userID, filter)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar penarikan berhasil diambil", result)
}

// GetWithdrawalByID godoc
// @Summary      Detail Withdraw
// @Description  Mengambil detail satu permintaan withdraw milik partner (proteksi IDOR)
// @Tags         Partner Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Withdraw ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /partner/wallet/withdrawals/{id} [get]
func (h *PartnerWalletHandler) GetWithdrawalByID(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "User tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	withdrawID := c.Param("id")
	result, err := h.usecase.GetWithdrawalByID(c.Request.Context(), userID, withdrawID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail penarikan berhasil diambil", result)
}

// RegisterRoutes mendaftarkan semua endpoint wallet partner
func (h *PartnerWalletHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/wallet", h.GetWalletDashboard)
	router.GET("/wallet/transactions", h.GetTransactions)
	router.POST("/wallet/withdraw", h.CreateWithdraw)
	router.GET("/wallet/withdrawals", h.GetMyWithdrawals)
	router.GET("/wallet/withdrawals/:id", h.GetWithdrawalByID)
}
