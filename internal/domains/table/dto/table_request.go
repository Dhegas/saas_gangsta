package dto

// CreateTableRequest payload untuk POST /api/dining-tables
type CreateTableRequest struct {
	TableName string `json:"table_name" binding:"required,max=50"`
}

// UpdateTableRequest payload untuk PUT /api/dining-tables/:id
type UpdateTableRequest struct {
	TableName string `json:"table_name" binding:"required,max=50"`
}
