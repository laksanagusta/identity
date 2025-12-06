package dtos

import "github.com/invopop/validation"

type ChangePassword struct {
	UserUUID    string `params:"userUUID"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (r ChangePassword) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.NewPassword, validation.Required, validation.Length(8, 255)),
		validation.Field(&r.OldPassword, validation.Required, validation.Length(8, 255)),
	)
}
