package main

import "errors"

var (
	ErrorPermNotInitialized = errors.New("permutation were not initialized by user")
	ErrorInvalidInput = errors.New("invalid input")
	ErrorDataReading = errors.New("data reading error")
)

type ErrorMsg struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func NewErrorMsg(code int, err error) ErrorMsg {
	return ErrorMsg{
		Code: code,
		Message: err.Error(),
	}
}