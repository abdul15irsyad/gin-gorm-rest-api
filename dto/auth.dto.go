package dto

type LoginDto struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterDto struct {
	Name            string `json:"name" form:"name" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
}

type UpdateAuthUserDto struct {
	Name  string `json:"name" form:"name" validate:"required"`
	Email string `json:"email" form:"email" validate:"required,email"`
}

type UpdateAuthUserPasswordDto struct {
	Password        string `json:"password" form:"password" validate:"required"`
	NewPassword     string `json:"newPassword" form:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required_with=NewPassword,eqfield=NewPassword"`
}
