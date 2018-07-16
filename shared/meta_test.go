package shared

import (
	"fmt"
	"testing"
	"time"
)

func TestAPIMetadata(t *testing.T) {
	expected := fmt.Sprintf("RadioChecker API (C) %d The RadioChecker Authors. "+
		"All rights reserved.", time.Now().Year())
	got := APIMetadata()
	if got != expected {
		t.Errorf("TestAPIMetadata: expected `%s`, got `%s`", expected, got)
	}
}
