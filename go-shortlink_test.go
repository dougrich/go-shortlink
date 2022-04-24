package shrtlink_test

import (
	"testing"

	"github.com/dougrich/go-shortlink"
)

func TestPlaceholder(t *testing.T) {
	if shrtlink.Placeholder() != 2 {
		t.Fatal("Expected a placeholder value of 2")
	}
}
