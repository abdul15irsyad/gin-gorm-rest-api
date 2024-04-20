package dto

type CreateStudentDto struct {
	Name string `json:"name" validate:"required"`
	Year int64  `json:"year" validate:"required,number,min=0"`
}

type UpdateStudentDto struct {
	Name string `json:"name" validate:"required,lowercase"`
	Year int64  `json:"year" validate:"required,numeric"`
}
