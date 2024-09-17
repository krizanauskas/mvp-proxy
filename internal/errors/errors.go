package errors

import "fmt"

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
