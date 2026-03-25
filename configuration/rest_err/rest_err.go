package rest_err

import "net/http"

type RestErr struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     string `json:"error"`
	Causes  []Causes
}

type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewRestError(message string, code int, err string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Code:    code,
		Err:     err,
		Causes:  causes,
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     "bad_request",
	}
}

func NewBadRequestValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     "internal_server_error",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Code:    http.StatusNotFound,
		Err:     "not_found",
	}
}
