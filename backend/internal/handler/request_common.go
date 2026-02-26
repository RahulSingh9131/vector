package handler

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// EmptyRequest is a no-op request for endpoints that don't accept input
type EmptyRequest struct{}

func (r *EmptyRequest) Validate() error {
	return nil
}

// EmptyResponse is a no-op response for endpoints that don't return data
type EmptyResponse struct{}
