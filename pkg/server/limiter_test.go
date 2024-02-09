package server_test

import (
	"github.com/jonboulle/clockwork"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIpRateLimiter(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()

	config := server.Config{Within: 10, RequestsLimit: 1, IpBucketLifeTimeSeconds: 30}
	rateLimiter := server.NewRateLimit(config, fakeClock)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(rateLimiter.IpRateLimiter(handler))
	defer srv.Close()

	// The first request must be http.StatusOk
	resp, err := http.Get(srv.URL)

	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("expected to get Status OK, got: %d", resp.StatusCode)
	}

	// The second request must be http.StatusTooManyRequests
	resp, err = http.Get(srv.URL)
	if err != nil || resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected to Status TooManyRequests, got: %d", resp.StatusCode)
	}

	// After 11 seconds in the future
	fakeClock.Advance(time.Second * 11)

	// The third request, after 2 seconds, must be http.StatusOk
	resp, err = http.Get(srv.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("expected third request to pass after time advance, got error: %v, status code: %d", err, resp.StatusCode)
	}
}
