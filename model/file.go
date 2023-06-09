package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type File struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Location string `json:"location"`
	FileType string `json:"fileType"`
	gorm.Model
}

type Files struct {
	Files []File
}

func (f *File) Set(filename string, location string, fileType string) {
	*f = File{
		Id:       filename + "+" + uuid.NewV4().String(),
		Filename: filename,
		Location: location,
		FileType: fileType,
	}
}
