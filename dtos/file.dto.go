package dtos

type GetFileDto struct {
	Id string `validate:"required,uuid"`
}
