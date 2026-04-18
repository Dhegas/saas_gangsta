package domain

type User struct {
	ID           string
	TenantID     string
	Email        string
	FullName     string
	PasswordHash string
	Role         string
	IsActive     bool
	TenantStatus string
}
