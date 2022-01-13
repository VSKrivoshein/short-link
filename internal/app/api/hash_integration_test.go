package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/api"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestRedirect(t *testing.T) {
	// Create user for test
	var (
		link     = "https://ya.ru/"
		email    = "test@email.com"
		password = "qwerty"
	)

	respSignUp := SignUp(t, email, password)
	defer RespClose(t, respSignUp)

	if respSignUp.StatusCode != http.StatusOK {
		t.Fatalf("error respSignUp.StatusCode != http.StatusOK")
	}

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password) //nolint:bodyclose
		defer RespClose(t, respDel)
		assert.Equal(t, http.StatusOK, respDel.StatusCode)
	}()

	respCreate := CreateLink(t, respSignUp, link) //nolint:bodyclose
	defer RespClose(t, respCreate)
	if respCreate.StatusCode != http.StatusOK {
		t.Fatalf("Fatal CreateLink(t, respSignUp, link): %v", respCreate.StatusCode)
	}
	out := new(api.CreateLinkResp)
	if err := json.NewDecoder(respCreate.Body).Decode(out); err != nil {
		t.Fatalf("Fatal json.NewDecoder(respCreate.Body).Decode(out): %v", err)
	}

	u, err := url.Parse(out.RedirectURL)
	if err != nil {
		t.Fatalf("Fatal url.Parse(out.RedirectURL): %v", err)
	}

	urlString := fmt.Sprintf("%v%v", TestSrv.URL, u.Path)
	req, err := http.NewRequest(http.MethodGet, urlString, http.NoBody) // nolint
	if err != nil {
		t.Fatalf("Fatal http.NewRequest(\"GET\", out.Link, nil): %v", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req) //nolint:bodyclose
	defer RespClose(t, resp)
	if err != nil {
		t.Fatalf("Fatal TestSrv.Client().Do(req): %v", err)
	}

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
}
