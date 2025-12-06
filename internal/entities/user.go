package entities

import (
	"time"

	"github.com/laksanagusta/identity/pkg/nullable"
)

type User struct {
	SoftDeleteModel
	EmployeeID       nullable.NullString `json:"employee_id" db:"employee_id"`
	Username         nullable.NullString `json:"username" db:"username"`
	Email            nullable.NullString `json:"email" db:"email"`
	FirstName        nullable.NullString `json:"first_name" db:"first_name"`
	LastName         nullable.NullString `json:"last_name" db:"last_name"`
	PhoneNumber      nullable.NullString `json:"phone_number" db:"phone_number"`
	PasswordHash     nullable.NullString `json:"-" db:"password_hash"`
	OrganizationUUID nullable.NullString `json:"organization_id" db:"organization_uuid"`
	IsActive         bool                `json:"is_active" db:"is_active"`
	LastLoginAt      *time.Time          `json:"last_login_at" db:"last_login_at"`

	Organization *Organization `json:"organization,omitempty" db:"-"`
	Roles        []*Role       `json:"roles,omitempty" db:"-"`
	Permissions  []*Permission `json:"permissions"`
}

type Users []*User

func (u *User) GetFullName() string {
	if u.LastName.IsExists {
		return u.FirstName.GetOrDefault() + " " + u.LastName.GetOrDefault()
	}
	return u.FirstName.GetOrDefault()
}

func (us Users) Uuids() []string {
	uuids := make([]string, 0, len(us))
	for _, u := range us {
		if u != nil {
			uuids = append(uuids, u.UUID)
		}
	}
	return uuids
}

type UserRole struct {
	BaseModel
	UserUUID string `json:"user_id" db:"user_uuid"`
	RoleUUID string `json:"role_id" db:"role_uuid"`
	User     *User  `json:"user,omitempty" db:"-"`
	Role     *Role  `json:"role,omitempty" db:"-"`
}

type ListUserParams struct {
	Offset    int
	Limit     int
	RoleUUID  nullable.NullString
	Role      nullable.NullString
	StartTime nullable.NullTime
	EndTime   nullable.NullTime
	Search    nullable.NullString
	Sort      *Sort
}
