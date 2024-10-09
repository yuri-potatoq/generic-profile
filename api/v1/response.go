package v1

import (
	"encoding/json"
	"fmt"
	"github.com/yuri-potatoq/generic-profile/pkg/pretty_errors"
	"net/http"
)

var mapErrors = map[pretty_errors.ErrorLevel]int{
	pretty_errors.UnrecoverableError:  http.StatusInternalServerError,
	pretty_errors.UnexpectedOperation: http.StatusBadRequest,
	pretty_errors.NoResources:         http.StatusNotFound,
	pretty_errors.NoAccess:            http.StatusUnauthorized,
	pretty_errors.OutOfBounds:         http.StatusForbidden,
}

type ErrorResponse struct {
	Message string          `json:"message"`
	Details map[string]any  `json:"details"`
	Errors  []ErrorResponse `json:"errors"`
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	var perr pretty_errors.PrettyError
	perr, ok := err.(pretty_errors.PrettyError)
	if !ok {
		WriteResponse(w, ErrorResponse{
			Message: fmt.Sprintf("unknown error: %s", err),
		}, http.StatusInternalServerError)
		return
	}
	//TODO: log here error details
	code, ok := mapErrors[perr.Level]
	if !ok {
		//TODO: log here
		code = http.StatusInternalServerError
	}
	WriteResponse(w, newErrorResponse(perr), code)
}

func newErrorResponse(perr pretty_errors.PrettyError) ErrorResponse {
	err := ErrorResponse{
		Message: perr.Message,
		Details: perr.Details,
	}

	for _, er := range perr.Inners() {
		err.Errors = append(err.Errors, newErrorResponse(er))
	}
	return err
}

func WriteResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		WriteErrorResponse(w, err)
		return
	}
}
