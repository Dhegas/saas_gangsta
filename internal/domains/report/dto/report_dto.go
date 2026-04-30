package dto

import "time"

// RevenueFilterParams parameter kueri untuk GET /api/reports/revenue
type RevenueFilterParams struct {
	From string `form:"from" binding:"required"` // Format: YYYY-MM-DD
	To   string `form:"to" binding:"required"`   // Format: YYYY-MM-DD
}

// TopMenusFilterParams parameter kueri untuk GET /api/reports/top-menus
type TopMenusFilterParams struct {
	From  string `form:"from" binding:"omitempty"`
	To    string `form:"to" binding:"omitempty"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// OrdersByTableFilterParams parameter kueri untuk GET /api/reports/orders-by-table
type OrdersByTableFilterParams struct {
	From  string `form:"from" binding:"omitempty"`
	To    string `form:"to" binding:"omitempty"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// DailySummaryFilterParams parameter kueri untuk GET /api/reports/daily-summary
type DailySummaryFilterParams struct {
	From string `form:"from" binding:"omitempty"` // Default: 7 hari ke belakang
	To   string `form:"to" binding:"omitempty"`   // Default: hari ini
}

// ---- Response Types ----

// RevenueResponse total pendapatan berdasarkan date range
type RevenueResponse struct {
	From         string    `json:"from"`
	To           string    `json:"to"`
	TotalRevenue float64   `json:"total_revenue"`
	TotalOrders  int       `json:"total_orders"`
	GeneratedAt  time.Time `json:"generated_at"`
}

// TopMenuEntry satu baris data menu terlaris
type TopMenuEntry struct {
	Rank      int     `json:"rank"`
	MenuID    string  `json:"menu_id"`
	MenuName  string  `json:"menu_name"`
	TotalQty  int     `json:"total_qty"`
	TotalSold float64 `json:"total_sold"`
}

// TopMenusResponse daftar menu terlaris
type TopMenusResponse struct {
	From  string         `json:"from,omitempty"`
	To    string         `json:"to,omitempty"`
	Menus []TopMenuEntry `json:"menus"`
}

// OrdersByTableEntry satu baris data order per meja
type OrdersByTableEntry struct {
	Rank        int     `json:"rank"`
	TableID     string  `json:"table_id"`
	TableNumber string  `json:"table_number"`
	TotalOrders int     `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
}

// OrdersByTableResponse daftar meja dengan jumlah order terbanyak
type OrdersByTableResponse struct {
	From   string               `json:"from,omitempty"`
	To     string               `json:"to,omitempty"`
	Tables []OrdersByTableEntry `json:"tables"`
}

// DailySummaryEntry ringkasan satu hari
type DailySummaryEntry struct {
	Date        string  `json:"date"`
	TotalOrders int     `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
	AvgOrderValue float64 `json:"avg_order_value"`
}

// DailySummaryResponse ringkasan harian dalam range tertentu
type DailySummaryResponse struct {
	From    string              `json:"from"`
	To      string              `json:"to"`
	Summary []DailySummaryEntry `json:"summary"`
}
