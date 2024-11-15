package dtos

type CreateUserDto struct {
	Name            string `json:"name" form:"name" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
	RoleId          string `json:"roleId" form:"roleId" validate:"uuid"`
}

type UpdateUserDto struct {
	Id              string  `json:"id" form:"id" validate:"required,uuid"`
	Name            string  `json:"name" form:"name" validate:"required"`
	Email           string  `json:"email" form:"email" validate:"required,email"`
	Password        *string `json:"password" form:"password" validate:"omitempty,min=8"`
	ConfirmPassword *string `json:"confirmPassword" form:"confirmPassword" validate:"required_with=Password,eqfield=Password"`
	RoleId          string  `json:"roleId" form:"roleId" validate:"uuid"`
}

type GetUserDto struct {
	Id string `validate:"required,uuid"`
}
