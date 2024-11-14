package models

type File struct {
	BaseModel
	Path             string `json:"path" gorm:"not null"`
	Filename         string `json:"filename" gorm:"not null"`
	OriginalFilename string `json:"originalFilename" gorm:"not null"`
	Mime             string `json:"mime" gorm:"not null"`
	Url              string `json:"url" gorm:"-"`
}

func (file *File) AfterLoad() {
	if file != nil {
		file.Url = "/assets" + file.Path + "/" + file.Filename
	}
}
