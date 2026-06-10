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

// AdminWalletHandler handler untuk endpoint wallet admin
type AdminWalletHandler struct {
	usecase walletdomain.AdminWalletUsecase
}

func NewAdminWalletHandler(usecase walletdomain.AdminWalletUsecase) *AdminWalletHandler {
	return &AdminWalletHandler{usecase: usecase}
}

// GetAllWithdrawals godoc
// @Summary      Daftar Semua Withdraw (Admin)
// @Description  Admin melihat semua permintaan withdraw dengan filter status opsional
// @Tags         Admin Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        status  query  string  false  "Filter: PENDING, APPROVED, TRANSFERRED, REJECTED"
// @Param        page    query  int     false  "Halaman (default: 1)"
// @Param        limit   query  int     false  "Jumlah per halaman (default: 20)"
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/wallet/withdrawals [get]
func (h *AdminWalletHandler) GetAllWithdrawals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	filter := walletdomain.WithdrawFilter{Status: status, Page: page, Limit: limit}
	result, err := h.usecase.GetAllWithdrawals(c.Request.Context(), filter)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar penarikan berhasil diambil", result)
}

// ApproveWithdrawal godoc
// @Summary      Setujui Withdraw
// @Description  Admin menyetujui permintaan withdraw partner (status: PENDING → APPROVED)
// @Tags         Admin Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Withdraw ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /admin/wallet/withdrawals/{id}/approve [patch]
func (h *AdminWalletHandler) ApproveWithdrawal(c *gin.Context) {
	adminID := c.GetString("userId")
	if adminID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Admin tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	withdrawID := c.Param("id")
	if err := h.usecase.ApproveWithdrawal(c.Request.Context(), adminID, withdrawID); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Withdraw berhasil disetujui", nil)
}

// RejectWithdrawal godoc
// @Summary      Tolak Withdraw
// @Description  Admin menolak permintaan withdraw (status: PENDING → REJECTED). Saldo otomatis dikembalikan ke wallet partner.
// @Tags         Admin Wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  string                     true  "Withdraw ID"
// @Param        request  body  dto.AdminReviewRequest  true  "Alasan penolakan"
// @Success      200  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /admin/wallet/withdrawals/{id}/reject [patch]
func (h *AdminWalletHandler) RejectWithdrawal(c *gin.Context) {
	adminID := c.GetString("userId")
	if adminID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Admin tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	withdrawID := c.Param("id")

	var req dto.AdminReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "admin_note diperlukan (min 3 karakter)", gin.H{
			"code": "VALIDATION_ERROR",
		})
		return
	}

	if err := h.usecase.RejectWithdrawal(c.Request.Context(), adminID, withdrawID, req.AdminNote); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Withdraw berhasil ditolak dan saldo dikembalikan", nil)
}

// MarkTransferred godoc
// @Summary      Tandai Transfer Selesai
// @Description  Admin menandai bahwa transfer manual ke rekening partner sudah dilakukan (status: PENDING/APPROVED → TRANSFERRED). total_withdrawn partner bertambah.
// @Tags         Admin Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Withdraw ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /admin/wallet/withdrawals/{id}/transfer [patch]
func (h *AdminWalletHandler) MarkTransferred(c *gin.Context) {
	adminID := c.GetString("userId")
	if adminID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Admin tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	withdrawID := c.Param("id")
	if err := h.usecase.MarkTransferred(c.Request.Context(), adminID, withdrawID); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Withdraw ditandai sebagai telah ditransfer", nil)
}

// GetAllPartnerWallets godoc
// @Summary      Daftar Wallet Semua Partner (Admin)
// @Description  Admin melihat ringkasan wallet seluruh partner di platform
// @Tags         Admin Wallet
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int  false  "Halaman (default: 1)"
// @Param        limit  query  int  false  "Jumlah per halaman (default: 20)"
// @Success      200  {object}  map[string]interface{}
// @Router       /admin/wallet/partners [get]
func (h *AdminWalletHandler) GetAllPartnerWallets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.usecase.GetAllPartnerWallets(c.Request.Context(), page, limit)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar wallet partner berhasil diambil", result)
}

// RegisterRoutes mendaftarkan semua endpoint wallet admin
func (h *AdminWalletHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/wallet/withdrawals", h.GetAllWithdrawals)
	router.PATCH("/wallet/withdrawals/:id/approve", h.ApproveWithdrawal)
	router.PATCH("/wallet/withdrawals/:id/reject", h.RejectWithdrawal)
	router.PATCH("/wallet/withdrawals/:id/transfer", h.MarkTransferred)
	router.GET("/wallet/partners", h.GetAllPartnerWallets)
}
