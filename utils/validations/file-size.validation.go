package validations

import (
	"mime/multipart"
)

type MaxFileSizeOptions struct {
	File    *multipart.FileHeader
	MaxSize *int64
}

var DefaultMaxSize int64 = 1024 * 1024 * 2

func MaxFileSize(options MaxFileSizeOptions) (ok bool) {
	if options.MaxSize == nil {
		options.MaxSize = &DefaultMaxSize
	}
	return options.File.Size <= *options.MaxSize
}
