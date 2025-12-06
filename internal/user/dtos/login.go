package dtos

import "github.com/invopop/validation"

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

func (r LoginReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Password, validation.Required, validation.Length(1, 255)),
	)
}
