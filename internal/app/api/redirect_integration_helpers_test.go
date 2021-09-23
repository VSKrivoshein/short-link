package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/api"
	"net/http"
	"testing"
)

func CreateLink(t *testing.T, response *http.Response, link string) *http.Response {
	url := fmt.Sprintf("%v/links/create", TestSrv.URL)

	content, err := json.Marshal(api.CreateLinkInput{Link: link})
	if err != nil {
		t.Fatalf("CreateLink Marshal Fatal: %v", err)
	}
	body := bytes.NewBuffer(content)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatalf("CreateLink http.Post Fatal: %v", err)
	}

	for _, cookie := range response.Cookies() {
		req.AddCookie(cookie)
	}

	resp, err := TestSrv.Client().Do(req)
	if err != nil {
		t.Fatalf("TestCreateLink Fatal TestSrv.Client().Do(req): %v", err)
	}

	return resp
}
