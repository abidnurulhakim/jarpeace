package response

import (
	"encoding/json"
	"net/http"
)

var (
	ErrTetapTenangTetapSemangat = CustomError{
		Message:  "Tetap Tenang Tetap Semangat",
		Code:     999,
		HTTPCode: http.StatusInternalServerError,
	}
)

type SuccessBody struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

type ErrorBody struct {
	Errors []ErrorInfo `json:"errors"`
	Meta   interface{} `json:"meta"`
}

type MetaInfo struct {
	HTTPStatus int `json:"http_status"`
}

type ErrorInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Field   string `json:"field,omitempty"`
}

type CustomError struct {
	Message  string
	Field    string
	Code     int
	HTTPCode int
}

func (c CustomError) Error() string {
	return c.Message
}

func BuildSuccess(data, meta interface{}) SuccessBody {
	return SuccessBody{
		Data: data,
		Meta: meta,
	}
}

func BuildError(errors []error) ErrorBody {
	var (
		ce CustomError
		ok bool
	)

	if len(errors) == 0 {
		ce = ErrTetapTenangTetapSemangat
	} else {
		err := errors[0]
		ce, ok = err.(CustomError)
		if !ok {
			ce = ErrTetapTenangTetapSemangat
		}
	}

	return ErrorBody{
		Errors: []ErrorInfo{
			{
				Message: ce.Message,
				Code:    ce.Code,
				Field:   ce.Field,
			},
		},
		Meta: MetaInfo{
			HTTPStatus: ce.HTTPCode,
		},
	}
}

func Write(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func WriteJson(w http.ResponseWriter, result interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}
