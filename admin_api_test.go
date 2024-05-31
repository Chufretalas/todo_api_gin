package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoginSendJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestAdmin_NoCredencials(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/admin/api/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Logf("expected responde code to be 401, got %v", w.Code)
		t.Fail()
	}
}

func TestAdmin_WrongCredentials(t *testing.T) {
	r := setupRouter()

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
	r := setupRouter()

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
