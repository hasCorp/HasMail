package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAuthToken(t *testing.T) {
	headers := make(http.Header)
	headers["Authorization"] = []string{"foo"}
	r := &http.Request{
		Header: headers,
	}

	assert.Equal(t, "foo", readAuthToken((r)))
}

func TestReadAuthTokenNoHeader(t *testing.T) {
	r := &http.Request{}
	assert.Empty(t, readAuthToken(r))
}

func TestReadAuthTokenEmptyHeader(t *testing.T) {
	headers := make(http.Header)
	headers["Authorization"] = []string{}
	r := &http.Request{
		Header: headers,
	}

	assert.Empty(t, readAuthToken(r))
}

func TestLocalAuthVerify(t *testing.T) {
	// token should be case sensitive
	tokens := []string{"foo", "Foo"}
	for _, token := range tokens {
		local := LocalAuthMiddleware{
			AllowedToken: token,
		}
		for _, compare := range tokens {
			r := &http.Request{
				Header: map[string][]string{
					"Authorization": {compare},
				},
			}
			t.Run(fmt.Sprintf("Compare %s to %s", token, compare), func(t *testing.T) {
				assert.Equal(t, token == compare, local.compareToLocalToken(r))
			})
		}
	}
}
