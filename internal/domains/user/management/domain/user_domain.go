package domain

import "time"

type User struct {
	ID        string
	TenantID  string
	Email     string
	FullName  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
