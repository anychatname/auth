package handlers

import (
	"fmt"
	"net/http/httptest"
	"strings"

	// "net/http/httptest"

	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	rR               RequestReaderImpl = RequestReaderImpl{}
	jsonNameExpected                   = "example"
	jsonOutput                         = struct {
		Name string `json:"name"`
	}{}
)

func TestRequestReaderImpl_JSON(t *testing.T) {
	t.Parallel()

	t.Run("Given valid both request and target output When reading request for JSON Then success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{ "name": "%s" }`, jsonNameExpected)))
		req.Header.Add("Content-Type", "application/json")

		err := rR.JSON(req, &jsonOutput)
		assert.NoError(t, err)
		assert.Equal(t, jsonNameExpected, jsonOutput.Name)
	})
}

func TestErrorRequestReaderImpl_JSON(t *testing.T) {
	t.Parallel()

	jsonOutput := struct {
		Name string `json:"name"`
	}{}

	t.Run("Given nil r param When checking r value Then invalid r *http.Request error", func(t *testing.T) {
		t.Parallel()

		err := rR.JSON(nil, &jsonOutput)
		assert.EqualError(t, err, "invalid request value: empty or nil *http.Request")
	})
	t.Run("Given nil r param When checking r value Then invalid r *http.Request error", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{ "name": "%s" }`, jsonNameExpected)))
		req.Header.Add("Content-Type", "application/json")

		err := rR.JSON(req, nil)
		assert.EqualError(t, err, "invalid target: empty or nil target object")
	})
	t.Run("Given invalid request content type When validating header content type Then invalid content type error", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{ "name": "%s" }`, jsonNameExpected)))
		err := rR.JSON(req, &jsonOutput)
		assert.EqualError(t, err, "invalid content type: Content-Type header is not application/json")
	})
	t.Run("Given empty JSON body When decoding JSON body from request Then decoding EOF error", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Add("Content-Type", "application/json")

		err := rR.JSON(req, &jsonOutput)
		assert.Contains(t, err.Error(), "error checking body: empty body content")
	})
	t.Run("Given invalid JSON body When decoding JSON body from request Then decoding error", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("POST", "/", strings.NewReader(`invalid JSON format`))
		req.Header.Add("Content-Type", "application/json")

		err := rR.JSON(req, &jsonOutput)
		assert.Contains(t, err.Error(), "error decoding body:")
	})
}
