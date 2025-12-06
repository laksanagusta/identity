package entities

import "github.com/google/uuid"

// AuthRole represents a user role from identity service (simplified for auth context)
type AuthRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserOrganization represents user organization from identity service
type UserOrganization struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

type AuthenticatedUser struct {
	ID           string           `json:"id"`
	EmployeeID   string           `json:"employee_id"`
	Username     string           `json:"username"`
	FirstName    string           `json:"first_name"`
	LastName     string           `json:"last_name"`
	PhoneNumber  string           `json:"phone_number"`
	Roles        []AuthRole       `json:"roles"`
	Organization UserOrganization `json:"organization"`
}
