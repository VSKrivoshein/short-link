package api_test

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/health", TestSrv.URL)) // nolint
	defer RespClose(t, resp)
	if err != nil {
		t.Fatalf("Expected no err")
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
