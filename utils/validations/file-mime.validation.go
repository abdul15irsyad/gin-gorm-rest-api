package validations

import (
	"gin-gorm-rest-api/utils"
	"mime/multipart"
)

type FileMimeOptions struct {
	File         *multipart.FileHeader
	AllowedMimes []string
}

func FileMime(options FileMimeOptions) (ok bool) {
	contentType := options.File.Header["Content-Type"][0]
	return utils.Contains(options.AllowedMimes, contentType)
}
