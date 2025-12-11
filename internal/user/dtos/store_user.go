package dtos

import (
	"regexp"

	"github.com/laksanagusta/identity/constants"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/errorhelper"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type CreateNewUserReq struct {
	EmployeeID       string              `json:"employee_id"`
	Username         string              `json:"username"`
	Password         string              `json:"password"`
	FirstName        string              `json:"first_name"`
	LastName         nullable.NullString `json:"last_name"`
	PhoneNumber      string              `json:"phone_number"`
	RoleUUIDs        []string            `json:"role_ids"`
	OrganizationUUID string              `json:"organization_id"`
}

// Password validation rules
var (
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasLowercase = regexp.MustCompile(`[a-z]`)
	hasNumber    = regexp.MustCompile(`[0-9]`)
	hasSymbol    = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`)
)

// ValidatePasswordStrength validates password strength
var ValidatePasswordStrength = validation.NewStringRuleWithError(
	func(s string) bool {
		return hasUppercase.MatchString(s) &&
			hasLowercase.MatchString(s) &&
			hasNumber.MatchString(s) &&
			hasSymbol.MatchString(s)
	},
	validation.NewError("validation_password_weak", "password harus mengandung kombinasi huruf besar, huruf kecil, angka, dan simbol"),
)

func (r CreateNewUserReq) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.EmployeeID, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Username, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 255), ValidatePasswordStrength),
		validation.Field(&r.FirstName, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.LastName, validation.Length(1, 255)),
		validation.Field(&r.PhoneNumber, validation.Required, is.UTFNumeric, validation.Length(11, 13)),
		validation.Field(&r.OrganizationUUID, validation.Required, is.UUIDv4),
		validation.Field(&r.RoleUUIDs, validation.Required, validation.Each(is.UUIDv4)),
	)

	roleUUIDSet := make(map[string]struct{}, len(r.RoleUUIDs))
	for _, uuid := range r.RoleUUIDs {
		if _, exists := roleUUIDSet[uuid]; exists {
			return errorhelper.BadRequestMap(map[string][]string{
				"role_ids": {constants.ErrDuplicated},
			})
		}
		roleUUIDSet[uuid] = struct{}{}
	}

	return err
}

func (r CreateNewUserReq) NewUser() entities.User {
	user := entities.User{
		EmployeeID:       nullable.NewString(r.EmployeeID),
		Username:         nullable.NewString(r.Username),
		FirstName:        nullable.NewString(r.FirstName),
		LastName:         r.LastName,
		PhoneNumber:      nullable.NewString(r.PhoneNumber),
		OrganizationUUID: nullable.NewString(r.OrganizationUUID),
	}

	for _, role := range r.RoleUUIDs {
		user.Roles = append(user.Roles, &entities.Role{
			SoftDeleteModel: entities.SoftDeleteModel{
				BaseModel: entities.BaseModel{
					UUID: role,
				},
			},
		})
	}

	user.BaseModel = entities.NewBaseModel("admin")

	return user
}
