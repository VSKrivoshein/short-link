package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/api"
	"io"
	"net/http"
	"testing"
)

var (
	contentType = "application/json"
)

func SignUp(t *testing.T, email, password string) *http.Response {
	url := fmt.Sprintf("%v/auth/sign-up", TestSrv.URL)
	content, err := json.Marshal(api.SignUpInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("CreateUser marshal fatal: %v", err)
	}

	body := bytes.NewBuffer(content)
	resp, err := http.Post(url, contentType, body) // nolint
	if err != nil {
		t.Fatalf("SignUp http.Post fatal: %v", err)
	}

	return resp
}

func DeleteUser(t *testing.T, email, password string) *http.Response {
	urlDelete := fmt.Sprintf("%v/auth/delete-user", TestSrv.URL)
	content, err := json.Marshal(api.SignUpInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Marshal Fatal: %v", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, urlDelete, bytes.NewBuffer(content)) // nolint
	if err != nil {
		t.Fatalf("Fatal http.NewRequest(\"DELETE\", contentType, bytes.NewBuffer(content)): %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Fatal resp, err := client.Do(req): %v", err)
	}

	return resp
}

func SignIn(t *testing.T, email, password string) *http.Response {
	url := fmt.Sprintf("%v/auth/sign-in", TestSrv.URL)
	content, err := json.Marshal(api.SignUpInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("CreateUser marshal fatal: %v", err)
	}
	body := bytes.NewBuffer(content)
	resp, err := http.Post(url, contentType, body) // nolint
	if err != nil {
		t.Fatalf("SignIn http.Post fatal: %v", err)
	}
	return resp
}

func RespToString(t *testing.T, resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("respToString: Fatal ioutil.ReadAll: %v", err)
	}
	return string(body)
}

func RespClose(t *testing.T, resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("Fatal closing resp")
	}
}
