package response

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExpectedSuccess struct {
	Test string
}

func TestSuccess(t *testing.T) {
	expectStatusCode := http.StatusOK
	expectedResult := ExpectedSuccess{Test: "test"}
	expectedBody := `{"test":"test"}`
	w := httptest.NewRecorder()
	err := NewSuccess(expectedResult, expectStatusCode)
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
