package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriterImpl_JSON(t *testing.T) {
	rW := ResponseWriterImpl{}

	t.Run("Given valid response writer When writing a JSON response Then success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		jsonExpected := jsonOutput
		jsonGot := jsonOutput
		jsonExpected.Name = jsonNameExpected
		assert.NoError(t, rW.JSON(rec, http.StatusOK, jsonExpected))

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, json.NewDecoder(rec.Body).Decode(&jsonGot))
		assert.Equal(t, jsonExpected, jsonGot)
	})
}
