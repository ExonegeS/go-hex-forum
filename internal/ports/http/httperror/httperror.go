package httperror

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-hex-forum/pkg/svcerr"
)

type APIError struct {
	Message    string
	StatusCode int
}

func FromError(err error) APIError {
	var svcErr *svcerr.Error
	apiErr := APIError{
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	}
	if errors.As(err, &svcErr) {
		apiErr.Message = svcErr.Message
		switch svcErr.AppErr {
		case svcerr.ErrNotAuthorized:
			apiErr.StatusCode = http.StatusUnauthorized
		case svcerr.ErrBadRequest:
			apiErr.StatusCode = http.StatusBadRequest
		case svcerr.ErrNotFound:
			apiErr.StatusCode = http.StatusNotFound
		case svcerr.ErrConflict:
			apiErr.StatusCode = http.StatusConflict
		case svcerr.ErrInternal:
			apiErr.StatusCode = http.StatusInternalServerError
		default:
			return APIError{"internal server error", http.StatusInternalServerError}
		}
	}

	return apiErr
}

func WriteError(w http.ResponseWriter, err error) {
	apiErr := FromError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": apiErr.Message,
	})
}
