package data

import "carbonite/admin/internal/validator"

type (
	LoginBody struct {
		Email    string `json:"email" form:"email" validate:"required,email"`
		Password string `json:"password" form:"password" validate:"required"`
	}
)

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
}

func ValidateLoginBody(v *validator.Validator, body *LoginBody) {
	ValidateEmail(v, body.Email)
	ValidatePasswordPlaintext(v, body.Password)
}
