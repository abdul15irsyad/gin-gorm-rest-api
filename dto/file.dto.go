package dto

type GetFileDto struct {
	Id string `validate:"required,uuid"`
}
