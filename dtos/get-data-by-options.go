package dtos

type GetDataByOptions struct {
	Field     string
	Value     string
	ExcludeId *string
}