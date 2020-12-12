package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func corsHandler(w http.ResponseWriter, r *http.Request) {
	if w.Header().Get(AccessExposeHeaders) != authorization ||
		w.Header().Get(AllowOrigin) != "*" ||
		w.Header().Get(AllowMethods) != AllowMethodsTypes ||
		w.Header().Get(AccessControlMaxAge) != "600" ||
		w.Header().Get(AllowHeaders) != ContentTypeeAuthorization {
		http.Error(w, "Cors header incomplete", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Cors header pass"))
}

func TestServeHTTP(t *testing.T) {
	msg := ""
	cases := []struct {
		caseName      string
		url           string
		expectedCode  int
		requestMethod string
		expectError   bool
	}{
		{
			"case session post bad requeest",
			"/v1/sessions",
			http.StatusBadRequest,
			http.MethodPost,
			true,
		},

		{
			"case session post must be in json",
			"/v1/sessions",
			http.StatusUnsupportedMediaType,
			http.MethodPost,
			true,
		},

		{
			"case cors",
			"/test/cors",
			http.StatusOK,
			http.MethodPost,
			false,
		},
		{
			"case test preflight",
			"/test/cors",
			http.StatusOK,
			http.MethodOptions,
			false,
		},
	}
	for _, c := range cases {
		if c.expectedCode == http.StatusBadRequest {
			msg = "{json:json}"
		}
		w := httptest.NewRecorder()
		r, err := http.NewRequest(c.requestMethod, c.url, strings.NewReader(msg))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if http.StatusBadRequest == c.expectedCode {
			r.Header.Set(headerContentType, contentTypeJSON)
		}
		cors := &CORS{http.HandlerFunc(corsHandler)}
		cors.ServeHTTP(w, r)
		if c.requestMethod != http.MethodOptions {
			msg := w.Body.String()
			if c.expectError && http.MethodPost == c.requestMethod {
				if w.Code != c.expectedCode {
					t.Errorf("Expect %v but got %v", c.expectedCode, w.Code)
				}
			} else {
				if msg != "Cors header pass" {
					t.Errorf("[Case %s] Expect %s but got %s", c.caseName, msg, w.Body.String())
				}
			}
		}
		if http.MethodOptions == c.requestMethod {
			if c.expectedCode != w.Code {
				t.Errorf("Expect %v but got %v", c.expectedCode, w.Code)
			}
		}
	}
}
