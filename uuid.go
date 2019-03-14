package goovn

import (
	"strings"

	"github.com/google/uuid"
)

func newUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return "row" + strings.Replace(id.String(), "-", "_", -1), nil
}
