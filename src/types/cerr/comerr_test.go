package cerr

import (
	"errors"
	"testing"
)

func TestLogError(t *testing.T) {
	LogError(errors.New("test log err"))
}
