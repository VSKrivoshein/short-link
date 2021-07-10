package apiserver

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIServer_handleHello(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "first test",
			path: "/hello",
			want: "Hello",
		},
		{
			name: "first test",
			path: "/hello",
			want: "Hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAPIServer(NewConfig())
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
			s.configureRouter()
			s.router.ServeHTTP(rec, req)
			assert.Equal(t, rec.Body.String(), "Hello")
		})
	}
}
