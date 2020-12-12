package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleWare_ServeHTTP(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapper := NewCORS(handler)
	req := httptest.NewRequest(http.MethodOptions, "/v1", nil)
	rr := httptest.NewRecorder()
	wrapper.ServeHTTP(rr, req)
	resp := rr.Result()

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("allow origin is not *")
	}

	if resp.Header.Get("Access-Control-Allow-Methods") != "GET, PUT, POST, PATCH, DELETE" {
		t.Errorf("allow methods is not allowed")
	}

	if resp.Header.Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("allow headers not Content-Type, Authorization")
	}

	if resp.Header.Get("Access-Control-Expose-Headers") != "Authorization" {
		t.Errorf("expose header not Authorization")
	}

	if resp.Header.Get("Access-Control-Max-Age") != "600" {
		t.Errorf("Max age not 600")
	}
}
