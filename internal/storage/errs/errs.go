package errs

import "errors"

var (
	ErrURLIsExist    = errors.New("URL already exists")
	ErrURLIsNotExist = errors.New("URL does not exist")
)
