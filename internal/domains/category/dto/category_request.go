package dto

// CreateCategoryRequest payload untuk POST /api/categories
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=120"`
	Description string `json:"description" binding:"omitempty"`
	SortOrder   int    `json:"sort_order" binding:"omitempty"`
}

// UpdateCategoryRequest payload untuk PUT /api/categories/:id
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,max=120"`
	Description string `json:"description" binding:"omitempty"`
}

// ToggleCategoryActiveRequest payload untuk PATCH /api/categories/:id/toggle-active
type ToggleCategoryActiveRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

// CategoryOrder item urutan
type CategoryOrder struct {
	ID        string `json:"id" binding:"required,uuid"`
	SortOrder int    `json:"sort_order" binding:"required"`
}

// ReorderCategoryRequest payload untuk PATCH /api/categories/reorder
type ReorderCategoryRequest struct {
	Items []CategoryOrder `json:"items" binding:"required,min=1"`
}
