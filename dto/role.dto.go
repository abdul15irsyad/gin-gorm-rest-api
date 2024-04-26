package dto

type GetRoleDto struct {
	Id string `validate:"required,uuid"`
}
