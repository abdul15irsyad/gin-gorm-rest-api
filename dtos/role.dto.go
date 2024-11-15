package dtos

type CreateRoleDto struct {
	Name string  `json:"name" form:"name" validate:"required"`
	Desc *string `json:"desc" form:"desc" validate:"omitempty"`
}

type UpdateRoleDto struct {
	Id   string  `json:"id" form:"id" validate:"required,uuid"`
	Name string  `json:"name" form:"name" validate:"required"`
	Desc *string `json:"desc" form:"desc" validate:"omitempty"`
}

type GetRoleDto struct {
	Id string `validate:"required,uuid"`
}
