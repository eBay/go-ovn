package goovn

import (
	"testing"

	"github.com/google/uuid"
)

func newUUID(t *testing.T) string {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}
