package errors

import (
	"errors"
	"net/http"

	"github.com/lib/pq"
)

const (
	DataNotFound          = "Data not found"
	DataAlreadyExist      = "Data already exist"
	ErrFeaturedSlotFull   = "featured slot is already occupied"
	ErrMaxFeaturedReached = "maximum 3 featured articles reached"
	ErrInvalidPosition    = "invalid featured position, must be 1, 2, or 3"
)

func CheckUniqueViolation(err error) error {
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	if ok && pqErr.Code == "23505" {
		return NewDefaultError(http.StatusBadRequest, DataAlreadyExist)
	}

	return err
}
