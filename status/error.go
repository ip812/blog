package status

import (
	"fmt"
	"net/http"
)

var (
	ErrDatabaseNotReady        = fmt.Errorf("database not initialized")
	ErrDB                      = fmt.Errorf("unexpected database error")
	ErrParsingFrom             = fmt.Errorf("failed to parse a form")
	ErrDecodingForm            = fmt.Errorf("failed to decode a form")
	ErrFailedtoValidateRequest = fmt.Errorf("failed to validate a request")

	ErrCreateArticleComment  = fmt.Errorf("failed to create an article comment")
	ErrGetAllArticleComments = fmt.Errorf("failed to get all article's comments")
)

func ErrorNotFound(err error) Toast {
	return Toast{
		Message:    err.Error(),
		StatusCode: http.StatusNotFound,
	}
}

func ErrorInternalServerError(err error) Toast {
	return Toast{
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}
