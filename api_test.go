package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

var r *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	r = setupRouter(false)
}

func TestAdmin_NoCredencials(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/admin/api/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Logf("expected responde code to be 401, got %v", w.Code)
		t.Fail()
	}
}

func TestAdmin_WrongCredentials(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/admin/api/users", nil)
	req.SetBasicAuth("aha", "aha")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Logf("expected responde code to be 401, got %v", w.Code)
		t.Fail()
	}
}

func TestAdmin_GetAllUsers(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/admin/api/users", nil)
	req.SetBasicAuth(adminUser, adminPass)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Logf("expected responde code to be 200, got %v", w.Code)
		t.Fail()
	}

	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Logf("expected responde Content-Type to be \"application/json; charset=utf-8\", got \"%v\"", w.Header().Get("Content-Type"))
		t.Fail()
	}
}

func TestAdmin_UserCreationAndDeletion(t *testing.T) {

	username := "test1"
	password := "test1"

	// making sure the user does not already exists
	w := httptest.NewRecorder()

	userJson, _ := json.Marshal(map[string]string{"username": username, "password": password})
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJson)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code == 200 {
		fmt.Printf("w.Code: %v\n", w.Code)
		token := gjson.Get(w.Body.String(), "token")
		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/validate", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		r.ServeHTTP(w, req)
		if w.Code == 200 {
			req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/api/users/%v", gjson.Get(w.Body.String(), "user.ID")), nil)
			req.SetBasicAuth(adminUser, adminPass)
			r.ServeHTTP(w, req)
		}
	}

	// Creating user
	w = httptest.NewRecorder()

	userJson, _ = json.Marshal(map[string]string{"username": username, "password": password})
	req = httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(userJson)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Logf("creating user /signup: expected responde code to be 200, got %v", w.Code)
		t.FailNow()
	}

	// Logging in
	w = httptest.NewRecorder()

	userJson, _ = json.Marshal(map[string]string{"username": username, "password": password})
	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJson)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Logf("logging user /login: expected responde code to be 200, got %v", w.Code)
		t.Fail()
	}

	token := gjson.Get(w.Body.String(), "token")
	if !token.Exists() {
		t.Log("logging user /login: response body did not contain \"token\" key")
		t.Fail()
	}

	// Validating
	w = httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/validate", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Logf("validating user /validate: expected responde code to be 200, got %v", w.Code)
		t.Fail()
	}

	userId := gjson.Get(w.Body.String(), "user.ID")
	if !userId.Exists() {
		t.Log("validating user /validate: response body did not contain \"ID\" key inside \"user\"")
		t.Fail()
	}

	// Deleting User
	w = httptest.NewRecorder()

	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/api/users/%v", userId), nil)
	req.SetBasicAuth(adminUser, adminPass)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Logf("deleting user /admin/api/users/%v: expected responde code to be 200, got %v", userId, w.Code)
		t.Fail()
	}

}
