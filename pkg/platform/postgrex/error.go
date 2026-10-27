package postgrex

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrQueryPreparation       = errors.New("query preparation error")
	ErrConcurrentModification = errors.New("concurrent modification detected")
)

func IsIntegrityViolation(err error) bool {
	var pgErr *pq.Error
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code.Class() == "23"
}
