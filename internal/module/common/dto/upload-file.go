package dto

import "mime/multipart"

type UploadFile struct {
	*multipart.FileHeader
	Module string `form:"module"`
}

func (f *UploadFile) SetDefault() {
	if f.Module == "" {
		f.Module = "common"
	}
}
