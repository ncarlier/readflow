package postgres

import (
	"strings"

	"github.com/ncarlier/readflow/internal/model"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key") {
		return model.ErrAlreadyExists
	}
	return err
}
