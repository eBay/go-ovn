package goovn

import (
	"github.com/google/uuid"
)

func newUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
