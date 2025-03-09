package common

import (
	"encoding/json"
	"net/http"

	"github.com/sant470/accounts-svc/common/errors"
)

type Handler func(rw http.ResponseWriter, r *http.Request) *errors.AppError

func (h Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := h(rw, r); err != nil {
		rw.WriteHeader(err.StatusCode)
		rw.Write([]byte(err.Message))
	}
}

func Decode(r *http.Request, v interface{}) *errors.AppError {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return errors.BadRequest(err.Error())
	}
	return nil
}
