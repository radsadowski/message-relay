package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	healthCh := make(chan struct{}, 1)
	handler := healthCheckHandler(healthCh)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler(res, req)

	assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)
	assert.Equal(t, struct{}{}, <-healthCh)
}
