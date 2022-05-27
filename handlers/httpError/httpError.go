package httpError

import (
	"net/http"
	"errors"
)

var (
	NotFoundError 		error = errors.New("404 Not Found")
	UnauthorizedError	error = errors.New("401 Unauthorised")
)

func GetStatusCode(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case NotFoundError:
		return http.StatusNotFound
	case UnauthorizedError:
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}