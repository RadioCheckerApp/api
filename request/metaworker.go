package request

import (
	"fmt"
	"time"
)

type MetaWorker struct{}

func (worker MetaWorker) HandleRequest() (interface{}, error) {
	formatStr := "RadioChecker API (C) %d The RadioChecker Authors. All rights reserved."
	return fmt.Sprintf(formatStr, time.Now().Year()), nil
}
