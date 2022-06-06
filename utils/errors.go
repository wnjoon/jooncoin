package utils

import "errors"

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

var ErrBlockNotFound = errors.New("block not found")
var ErrNotEnoughMoney = errors.New("not enough money")
var ErrNotValidated = errors.New("not validated")
