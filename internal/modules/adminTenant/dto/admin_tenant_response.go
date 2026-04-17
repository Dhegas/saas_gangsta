package dto

import "time"

// TenantResponse adalah data yang akan dikirim ke frontend Admin
type TenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"` // active, inactive, suspended
	CreatedAt time.Time `json:"createdAt"`
}
