package dto

// CreateMenuRequest payload untuk POST /api/menus
type CreateMenuRequest struct {
	CategoryID  *string  `json:"category_id" binding:"omitempty,uuid"`
	Name        string   `json:"name" binding:"required,max=180"`
	Description string   `json:"description" binding:"omitempty"`
	Price       float64  `json:"price" binding:"required,min=0"`
	ImageURL    string   `json:"image_url" binding:"omitempty,url"`
}

// UpdateMenuRequest payload untuk PUT /api/menus/:id
type UpdateMenuRequest struct {
	CategoryID  *string  `json:"category_id" binding:"omitempty,uuid"`
	Name        string   `json:"name" binding:"omitempty,max=180"`
	Description string   `json:"description" binding:"omitempty"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	ImageURL    string   `json:"image_url" binding:"omitempty,url"`
}

// ToggleMenuAvailableRequest payload untuk PATCH /api/menus/:id/toggle-available
type ToggleMenuAvailableRequest struct {
	IsAvailable *bool `json:"is_available" binding:"required"`
}

// MenuFilterParams parameter untuk GET /api/menus
type MenuFilterParams struct {
	CategoryID  string `form:"category_id" binding:"omitempty,uuid"`
	IsAvailable *bool  `form:"is_available" binding:"omitempty"`
}
