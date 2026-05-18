package dto

// CreateTableRequest payload untuk POST /api/v1/tables
type CreateTableRequest struct {
	TableName string `json:"table_name" binding:"required,max=50"`
}

// UpdateTableRequest payload untuk PUT /api/v1/tables/:id
type UpdateTableRequest struct {
	TableName string `json:"table_name" binding:"required,max=50"`
}
