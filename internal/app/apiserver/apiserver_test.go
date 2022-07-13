package apiserver

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_Start(t *testing.T) {
	config, err := NewConfigBuilder().Parse("../../../configs/apiserver.json")
	assert.Nil(t, err)
	s := NewAPIServer(config)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.handlePost().ServeHTTP(rec, req)

	assert.Equal(t, rec.Body.String(), "Hello")
}
