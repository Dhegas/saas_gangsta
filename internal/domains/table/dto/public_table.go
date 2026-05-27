package dto

type PublicTableResponse struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	TableName string `json:"tableName"`
	Status    string `json:"status"` // "occupied" or "kosong"
}
