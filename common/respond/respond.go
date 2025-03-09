package respond

import (
	"encoding/json"
	"net/http"

	"github.com/sant470/accounts-svc/common/errors"
)

func toJSON(rw http.ResponseWriter, status int, data interface{}) *errors.AppError {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return errors.InternalServerError(err.Error())
	}
	return nil
}

func OK(rw http.ResponseWriter, data interface{}) *errors.AppError {
	return toJSON(rw, http.StatusOK, data)
}
