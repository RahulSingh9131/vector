package errs

import "net/http"

func NewUnauthorizedError(message string, override bool) *HTTPError {
	return &HTTPError{
		Code:     MakeUpperCaseWithUnderscore(http.StatusText(http.StatusUnauthorized)),
		Message:  message,
		Status:   http.StatusUnauthorized,
		Override: override,
	}
}

func NewForbiddenError(message string, override bool) *HTTPError {
	return &HTTPError{
		Code:     MakeUpperCaseWithUnderscore(http.StatusText(http.StatusForbidden)),
		Message:  message,
		Status:   http.StatusForbidden,
		Override: override,
	}
}

func NewBadRequestError(message string, override bool, code *string, errors []FieldError, action *Action) *HTTPError {
	formattedCode := MakeUpperCaseWithUnderscore(http.StatusText(http.StatusBadRequest))

	if code != nil {
		formattedCode = *code
	}

	return &HTTPError{
		Code:     formattedCode,
		Message:  message,
		Status:   http.StatusBadRequest,
		Override: override,
		Errors:   errors,
		Action:   action,
	}
}

func NewNotFoundError(message string, override bool, code *string) *HTTPError {
	formattedCode := MakeUpperCaseWithUnderscore(http.StatusText(http.StatusNotFound))
	if code != nil {
		formattedCode = *code
	}
	return &HTTPError{
		Code:     formattedCode,
		Message:  message,
		Status:   http.StatusNotFound,
		Override: override,
	}
}

func NewInternalServerError(message string, override bool) *HTTPError {
	return &HTTPError{
		Code:     MakeUpperCaseWithUnderscore(http.StatusText(http.StatusInternalServerError)),
		Message:  message,
		Status:   http.StatusInternalServerError,
		Override: override,
	}
}

func ValidationError(error error) *HTTPError {
	return NewBadRequestError("validation failed"+error.Error(), false, nil, nil, nil)
}
