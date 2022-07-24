package sample

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hascorp/hasmail/internal/hasmailtemplates"
	"github.com/hascorp/hasmail/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// need to init env vars and mock API client
	testutil.MockEmailClient()

	os.Exit(m.Run())
}

func TestSendSampleEmail(t *testing.T) {
	ctx := context.TODO()

	body := hasmailtemplates.MailRequest{
		Recipient: "hank@hascorp.dev",
		Name:      "Hank Pecker",
		Vars: map[string]string{
			"foo": "PepeLa",
		},
	}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, "POST", "/mail/sample", bytes.NewBuffer(b))
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Need a new router in order to inject the mux vars correctly
	router := mux.NewRouter()
	router.HandleFunc("/mail/sample", SampleHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusAccepted, rr.Code)

	// Check the response body is what we expect.
	expected := `{"success":true}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestSendSampleEmailError(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testutil.FailWithCode{},
		http.StatusForbidden,
	)

	body := hasmailtemplates.MailRequest{
		Recipient: "hank@hascorp.dev",
		Name:      "Hank Pecker",
		Vars: map[string]string{
			"foo": "PepeLa",
		},
	}
	b, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, "POST", "/mail/sample", bytes.NewBuffer(b))
	assert.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Need a new router in order to inject the mux vars correctly
	router := mux.NewRouter()
	router.HandleFunc("/mail/sample", SampleHandler)
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Check the response body is what we expect.
	expected := `{"success":false}`
	assert.Equal(t, expected, rr.Body.String())
}
