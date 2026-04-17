package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Pastikan path import ini sesuai dengan nama modul di go.mod kamu
	tenantHttp "github.com/dhegas/saas_gangsta/internal/modules/adminTenant/delivery/http"
	tenantRepo "github.com/dhegas/saas_gangsta/internal/modules/adminTenant/repository"
	tenantUc "github.com/dhegas/saas_gangsta/internal/modules/adminTenant/usecase"
)

// RegisterAdminRoutes adalah titik masuk untuk mendaftarkan semua endpoint khusus Admin.
// Fungsi ini menerima grup router utama (misal: /api/v1) dan koneksi database.
func RegisterAdminRoutes(apiV1 *gin.RouterGroup, db *gorm.DB) {

	// ==========================================
	// 1. Modul Admin Tenant
	// ==========================================

	// Rangkai dependensi (Dependency Injection)
	adminTenantRepo := tenantRepo.NewAdminTenantRepository(db)
	adminTenantUsecase := tenantUc.NewAdminTenantUsecase(adminTenantRepo)
	adminTenantHandler := tenantHttp.NewTenantHandler(adminTenantUsecase)

	// Daftarkan rute ke dalam grup /api/v1
	// Catatan: Nanti kita bisa tambahkan middleware auth di sini
	adminTenantHandler.RegisterRoutes(apiV1)

	// ==========================================
	// 2. Modul Admin Dashboard (Ruang Kosong untuk Nanti)
	// ==========================================
	// Nanti kode untuk inisialisasi modul dashboard ditaruh di sini

	// ==========================================
	// 3. Modul Admin Subscription (Ruang Kosong untuk Nanti)
	// ==========================================
	// Nanti kode untuk inisialisasi modul subscription ditaruh di sini
}
