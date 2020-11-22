package errors

import (
	goerrors "errors"
	"testing"
)

func TestAnnotate(t *testing.T) {
	err := goerrors.New("abc123")
	newErr := Annotate(err, CodeInternal)

	got, want := newErr.Error(), err.Error()
	if got != want {
		t.Errorf("newErr.Error(): got %v, want %v", got, want)
	}
}
