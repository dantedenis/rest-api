package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_Start(t *testing.T) {
	config, err := NewConfigBuilder().Parse("../../../configs/apiserver.json")
	//assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
	s := NewAPIServer(config)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.handlePost().ServeHTTP(rec, req)

	if rec.Body.String() != "Hello" {
		t.Error("Invalid request")
	}

}


func BenchmarkAPIServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config, _ := NewConfigBuilder().Parse("../../../configs/apiserver.json")
		s := NewAPIServer(config)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		s.handlePost().ServeHTTP(rec, req)		
	}
}