package errors

import (
	"errors"
	"github.com/lib/pq"
	"net/http"
)

const (
	DataNotFound     = "Data not found"
	DataAlreadyExist = "Data already exist"
)

func CheckUniqueViolation(err error) error {
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	if ok && pqErr.Code == "23505" {
		return NewDefaultError(http.StatusBadRequest, DataAlreadyExist)
	}

	return err
}
