package noop

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hascorp/hasmail/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// need to init env vars and mock API client
	testutil.MockEmailClient()

	os.Exit(m.Run())
}

func TestSendNoOp(t *testing.T) {
	ctx := context.Background()

	body := strings.NewReader(`{"a":"b"}`)
	req, err := http.NewRequestWithContext(ctx, "POST", "/mail/noop", body)
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Need a new router in order to inject the mux vars correctly
	router := mux.NewRouter()
	router.HandleFunc("/mail/noop", NoOpHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusAccepted, rr.Code)

	// Check the response body is what we expect.
	expected := `{"success":true}`
	assert.Equal(t, expected, rr.Body.String())
}
