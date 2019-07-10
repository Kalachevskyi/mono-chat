package handlers

import "fmt"

// ErrStack - add stack to error, work with "github.com/pkg/errors" package
func ErrStack(err error) error {
	if err == nil {
		return err
	}

	return fmt.Errorf("%+v", err)
}
