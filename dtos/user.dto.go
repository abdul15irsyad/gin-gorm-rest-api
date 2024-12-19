package dtos

import (
	"mime/multipart"
)

type CreateUserDto struct {
	Name            string `json:"name" form:"name" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" validate:"required,eqfield=Password"`
	RoleId          string `json:"roleId" form:"roleId" validate:"uuid"`
}

type UpdateUserDto struct {
	Id              string                `json:"id" form:"id" validate:"required,uuid"`
	Name            string                `json:"name" form:"name" validate:"required"`
	Email           string                `json:"email" form:"email" validate:"required,email"`
	Password        *string               `json:"password" form:"password" validate:"omitempty,min=8"`
	ConfirmPassword *string               `json:"confirmPassword" form:"confirmPassword" validate:"required_with=Password,eqfield=Password"`
	RoleId          string                `json:"roleId" form:"roleId" validate:"uuid"`
	Image           *multipart.FileHeader `json:"image" form:"image"`
}

type GetUserDto struct {
	Id string `validate:"required,uuid"`
}

type GetUsersDto struct {
	Page   *int    `form:"page" validate:"omitempty,number,gte=1"`
	Limit  *int    `form:"limit" validate:"omitempty,number"`
	Search *string `form:"search" validate:"omitempty"`
}
