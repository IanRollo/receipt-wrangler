package middleware

import (
	"net/http/httptest"
	"testing"
)

// The processing queue routes, classified for API-key scopes. Kept beside the main
// inventory (everyRoute) so the fork routes do not sit in a table upstream edits.
// The three GETs disclose and the retry changes state without returning anything,
// so no allowlist entry is needed: the method decides, and fails closed.
var processingQueueRoutes = []routeCase{
	{"GET", "/api/processingQueue/group/1", classRead},
	{"GET", "/api/processingQueue/task/quick_scan/abc/sourceFile", classRead},
	{"GET", "/api/processingQueue/task/quick_scan/abc/sourceFile/download", classRead},
	{"POST", "/api/processingQueue/task/email_receipt_processing/abc/retry", classWrite},
}

func TestClassifyRequest_ProcessingQueueRoutes(t *testing.T) {
	for _, route := range processingQueueRoutes {
		r := httptest.NewRequest(route.method, route.path, nil)
		if got := classifyRequest(r); got != route.class {
			t.Errorf("%s %s: got %s, want %s", route.method, route.path, className(got), className(route.class))
		}
	}
}
