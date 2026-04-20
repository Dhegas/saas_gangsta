package dto

// DashboardStatsResponse berisi metrik utama untuk halaman depan Admin
type DashboardStatsResponse struct {
	TotalTenants        int     `json:"totalTenants"`
	ActiveSubscriptions int     `json:"activeSubscriptions"`
	MonthlyRevenue      float64 `json:"monthlyRevenue"`
}
