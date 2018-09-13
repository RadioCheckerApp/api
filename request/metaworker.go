package request

import (
	"fmt"
	"time"
)

type MetaWorker struct{}

var Version string
var Build string
var Revision string

func (worker MetaWorker) HandleRequest() (interface{}, error) {
	formatStr := "RadioChecker API (C) %d The RadioChecker Authors. All rights reserved. " +
		"(Version: %s / Build: %s / Revision: %s)"
	return fmt.Sprintf(formatStr, time.Now().Year(), Version, Build, Revision), nil
}
