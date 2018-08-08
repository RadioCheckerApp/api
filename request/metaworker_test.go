package request

import (
	"fmt"
	"testing"
	"time"
)

func TestMetaWorker_HandleRequest(t *testing.T) {
	formatStr := "RadioChecker API (C) %d The RadioChecker Authors. All rights reserved."
	expectedResult := fmt.Sprintf(formatStr, time.Now().Year())
	result, err := MetaWorker{}.HandleRequest()
	if err != nil || result != expectedResult {
		t.Errorf("MetaWorker (%v).HandleRequest(): got (%s, %v), expected type (%s, false)",
			MetaWorker{}, result, err, expectedResult)
	}
}
