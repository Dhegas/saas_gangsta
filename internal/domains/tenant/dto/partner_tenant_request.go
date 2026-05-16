package dto

import "mime/multipart"

type CreatePartnerTenantRequest struct {
	Name        string                `form:"name" binding:"required,min=2,max=150"`
	Description string                `form:"description"`
	Address     string                `form:"address"`
	PhoneNumber string                `form:"phone_number"`
	OpenHours   string                `form:"open_hours"`
	Logo        *multipart.FileHeader `form:"logo"`
	Banner      *multipart.FileHeader `form:"banner"`
}
