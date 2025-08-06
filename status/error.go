package status

import (
	"fmt"
	"net/http"
)

var (
	ErrDatabaseNotReady = fmt.Errorf("database not initialized")
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
