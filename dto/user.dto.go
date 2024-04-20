package dto

type CreateUserDto struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type UpdateUserDto struct {
	Name            string  `json:"name" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	Password        *string `json:"password" validate:"omitempty,min=8"`
	ConfirmPassword *string `json:"confirmPassword" validate:"required_with=Password,eqfield=Password"`
}
