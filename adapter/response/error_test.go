package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	expectStatusCode := http.StatusOK
	expecErr := errors.New("error test")
	expectedBody := `{"error":"` + expecErr.Error() + `"}`
	w := httptest.NewRecorder()
	err := NewError(expecErr, expectStatusCode, nil)
	err.Send(w)

	var result = strings.TrimSpace(w.Body.String())
	if !strings.EqualFold(result, expectedBody) {
		t.Errorf(
			"[TestCase '%s'] Result: '%v' | Expected: '%v'",
			"test",
			result,
			expectedBody,
		)
	}

	assert.Equal(t, w.Result().StatusCode, expectStatusCode)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
}
