package errors

import "fmt"

type InvalidField struct {
	Field string
	AsIs  string
}

var _ error = (*InvalidField)(nil)

func (e InvalidField) Error() string {
	return fmt.Sprintf("Field [%s] is invalid. %s.", e.Field, e.AsIs)
}
