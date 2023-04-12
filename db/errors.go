package db

import (
	"fmt"

	"github.com/lib/pq"
)

const (
	UniqueViolationError = pq.ErrorCode("23505")
)

func IsUniqueConstraintViolationError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == UniqueViolationError
	}

	return false
}

type BadInputError struct {
	Err error
}

func (b *BadInputError) Error() string {
	return fmt.Sprintf("%v", b.Err)
}

type NotFoundError struct {
	object string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("Could not find an %s with the provided parameters", n.object)
}
