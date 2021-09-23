package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/api"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateLink(t *testing.T) {
	// Create user for test
	var (
		link     = "https://ya.ru/"
		email    = "test@email.com"
		password = "qwerty"
	)

	respSignUp := SignUp(t, email, password)
	defer RespClose(t, respSignUp)

	if respSignUp.StatusCode != http.StatusOK {
		t.Fatalf("respSignUp.StatusCode != http.StatusOK")
	}

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, http.StatusOK, respDel.StatusCode)
	}()

	// Test creating link
	resp := CreateLink(t, respSignUp, link)
	defer RespClose(t, resp)

	respStruct := new(api.CreateLinkResp)
	if err := json.NewDecoder(resp.Body).Decode(respStruct); err != nil {
		t.Fatalf("TestCreateLink Fatal json.NewDecoder(resp.Body).Decode(respStruct) %v", err)
	}

	assert.Equal(t, link, respStruct.Link)
	if err := validator.New().Struct(respStruct); err != nil {
		t.Errorf("response redirect url is not valid url")
	}
}

func TestCreateExistedLink(t *testing.T) {
	// Create user for test
	var (
		link     = "https://ya.ru/"
		email    = "test@email.com"
		password = "qwerty"
	)

	respSignUp := SignUp(t, email, password)
	defer RespClose(t, respSignUp)

	if respSignUp.StatusCode != http.StatusOK {
		t.Fatalf("respSignUp.StatusCode != http.StatusOK")
	}

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, http.StatusOK, respDel.StatusCode)
	}()

	// Test creating link
	resp := CreateLink(t, respSignUp, link)
	defer RespClose(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("resp.StatusCode != http.StatusOK, code: %v", resp.StatusCode)
	}

	// Try Creating existed link
	respExisted := CreateLink(t, respSignUp, link)

	assert.NotEqual(t, http.StatusOK, respExisted.StatusCode)
}

func TestCreateLinkUnauthorized(t *testing.T) {
	// Create user for test
	var (
		link = "https://ya.ru/"
	)

	// Response without authorization cookies
	respSignUp := new(http.Response)

	// Test creating link
	resp := CreateLink(t, respSignUp, link)
	defer RespClose(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestCreateLinkInvalidUrl(t *testing.T) {
	// Create user for test
	var (
		link     = "invalid"
		email    = "test@email.com"
		password = "qwerty"
	)

	respSignUp := SignUp(t, email, password)
	defer RespClose(t, respSignUp)

	if respSignUp.StatusCode != http.StatusOK {
		t.Fatalf("respSignUp.StatusCode != http.StatusOK")
	}

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, http.StatusOK, respDel.StatusCode)
	}()

	// Test creating link
	resp := CreateLink(t, respSignUp, link)
	defer RespClose(t, resp)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestGetAllLinks(t *testing.T) {
	// Create user for test
	var (
		links    = []string{"https://ya.ru/", "https://mail.ru/", "https://vc.ru/"}
		email    = "test@email.com"
		password = "qwerty"
		url      = fmt.Sprintf("%v/links/get-all", TestSrv.URL)
	)

	respSignUp := SignUp(t, email, password)
	defer RespClose(t, respSignUp)

	if respSignUp.StatusCode != http.StatusOK {
		t.Fatalf("respSignUp.StatusCode != http.StatusOK")
	}

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, http.StatusOK, respDel.StatusCode)
	}()

	// Creating few links per user
	for _, link := range links {
		resp := CreateLink(t, respSignUp, link)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		RespClose(t, resp)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Fatal creating new request http.NewRequest: %v", err)
	}

	for _, cookie := range respSignUp.Cookies() {
		req.AddCookie(cookie)
	}

	resp, err := TestSrv.Client().Do(req)
	defer RespClose(t, resp)
	if err != nil {
		t.Fatalf("Fatal TestSrv.Client().Do(req): %v", err)
	}

	respStruct := new(api.GetAllLinksResp)
	if err := json.NewDecoder(resp.Body).Decode(respStruct); err != nil {
		t.Fatalf("Fatal json.NewDecoder(resp.Body).Decode(respStruct): %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, len(respStruct.AllUserLinks))
}

func TestDeleteLink(t *testing.T) {
	// Create user for test
	var (
		link     = "https://ya.ru/"
		email    = "test@email.com"
		password = "qwerty"
		url      = fmt.Sprintf("%v/links/delete", TestSrv.URL)
	)

	respSignUp := SignUp(t, email, password)
	defer RespClose(t, respSignUp)

	if respSignUp.StatusCode != http.StatusOK {
		t.Fatalf("respSignUp.StatusCode != http.StatusOK")
	}

	// rm user
	defer func() {
		respDel := DeleteUser(t, email, password)
		defer RespClose(t, respDel)
		assert.Equal(t, http.StatusOK, respDel.StatusCode)
	}()

	respCreate := CreateLink(t, respSignUp, link)
	defer RespClose(t, respCreate)
	if respCreate.StatusCode != http.StatusOK {
		t.Fatalf("Fatal CreateLink(t, respSignUp, link): %v", respCreate.StatusCode)
	}

	input := api.DeleteLinkInput{Link: link}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Fatal json.Marshal(input): %v", err)
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Fatal http.NewRequest: %v", err)
	}

	for _, cookie := range respSignUp.Cookies() {
		req.AddCookie(cookie)
	}

	resp, err := TestSrv.Client().Do(req)
	defer RespClose(t, resp)
	if err != nil {
		t.Fatalf("TestSrv.Client().Do(req): %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}