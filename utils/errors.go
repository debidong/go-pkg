package utils

import "fmt"

func NewError(errMsg, functionName string) error {
	return fmt.Errorf("%s: %s", functionName, errMsg)
}

func WrapError(err error, functionName string) error {
	return fmt.Errorf("%s: %w", functionName, err)
}
