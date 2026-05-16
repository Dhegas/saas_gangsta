package dto

import "time"

// TableResponse DTO response untuk daftar dan detail meja
type TableResponse struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	TableName string     `json:"table_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TableStatusResponse DTO response untuk status meja
type TableStatusResponse struct {
	ID        string `json:"id"`
	TableName string `json:"table_name"`
	Status    string `json:"status"` // "kosong" atau "occupied"
}
