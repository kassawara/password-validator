package errors

import "fmt"

type NotFoundError struct {
	Entity string
	ID     string
}

var _ error = (*NotFoundError)(nil)

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found with ID '%s'", e.Entity, e.ID)
}
