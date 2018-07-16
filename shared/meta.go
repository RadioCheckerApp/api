package shared

import (
	"fmt"
	"time"
)

var date = time.Now()

func APIMetadata() string {
	return fmt.Sprintf("RadioChecker API (C) %d The RadioChecker Authors. All rights reserved.",
		date.Year())
}
