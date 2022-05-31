package utils

import "errors"

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

var ErrBlockNotFound = errors.New("Block not found")
