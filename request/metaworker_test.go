package request

import (
	"fmt"
	"testing"
	"time"
)

func TestMetaWorker_HandleRequest(t *testing.T) {
	formatStr := "RadioChecker API (C) %d The RadioChecker Authors. All rights reserved. " +
		"(Version: %s / Build: %s / Revision: %s)"
	expectedResult := fmt.Sprintf(formatStr, time.Now().Year(), Version, Build, Revision)
	result, err := MetaWorker{}.HandleRequest()
	if err != nil || result != expectedResult {
		t.Errorf("MetaWorker (%v).HandleRequest(): got (%s, %v), expected type (%s, false)",
			MetaWorker{}, result, err, expectedResult)
	}
}
