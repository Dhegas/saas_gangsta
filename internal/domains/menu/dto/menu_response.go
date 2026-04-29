package dto

import "time"

// MenuResponse DTO response untuk daftar dan detail menu
type MenuResponse struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	CategoryID  *string    `json:"category_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	ImageURL    string     `json:"image_url"`
	IsAvailable bool       `json:"is_available"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
