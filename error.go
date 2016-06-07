package filepool

import (
	"errors"
)

var (
	ErrRequestDataType   error = errors.New("Request data field type error")
	ErrUploadFileSize    error = errors.New("Upload File size error")
	ErrRequestRestMethod error = errors.New("RESTful method error")
	ErrImageType         error = errors.New("Type of file is not image")
	ErrPermission        error = errors.New("error Permission")
)
