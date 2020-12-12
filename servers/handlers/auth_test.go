package handlers

import (
	"assignments-tichx/servers/gateway/models/users"
	"assignments-tichx/servers/gateway/sessions"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

func newUser() *users.NewUser {
	return &users.NewUser{
		UserName:     "info441",
		Email:        "testmail@testmail.com",
		FirstName:    "First",
		LastName:     "Last",
		Password:     "Password@123",
		PasswordConf: "Password@123",
	}
}
func newSpecificUser(id int64) *users.User {
	return &users.User{
		ID:        id,
		Email:     "testmail@testmail.com",
		PassHash:  []byte("Password@123"),
		UserName:  "username",
		FirstName: "first1",
		LastName:  "last1",
		PhotoURL:  "randomurl",
	}
}

func newEmptyUser() *users.User {
	return &users.User{}
}
func TestUsersHandler(t *testing.T) {
	nu := newUser()
	eu := newUser()
	eu.UserName = "username1"
	userStore := &users.FakeMySQLStore{}
	sessStore := sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), 3600)
	ctx := &HandlerContext{
		Key:          "somerandomkey",
		SessionStore: sessStore,
		UserStore:    userStore,
	}
	cases := []struct {
		caseName      string
		requestMethod string
		requestType   string
		expectedCode  int
		expectError   bool
		user          *users.NewUser
	}{
		{
			"Pass",
			http.MethodPost,
			contentTypeJSON,
			http.StatusCreated,
			false,
			nu,
		},

		{
			"Error: request content type incorrect",
			http.MethodPost,
			"text/plain",
			http.StatusUnsupportedMediaType,
			true,
			&users.NewUser{},
		},

		{
			"Error: error user insertion",
			http.MethodPost,
			contentTypeJSON,
			http.StatusInternalServerError,
			true,
			eu,
		},
		{
			"Error: empty user fail",
			http.MethodPost,
			contentTypeJSON,
			http.StatusBadRequest,
			true,
			&users.NewUser{},
		},
		{
			"Error: Method not allowed",
			http.MethodPatch,
			contentTypeJSON,
			http.StatusMethodNotAllowed,
			true,
			&users.NewUser{},
		},
	}
	for _, c := range cases {
		body, err := json.Marshal(c.user)
		if err != nil {
			t.Errorf("json encode went wrong")
		}
		handler := http.HandlerFunc(ctx.UsersHandler)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(c.requestMethod, "/v1/users", strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("$err")
		}
		r.Header.Set(headerContentType, c.requestType)
		if !c.expectError {
			handler.ServeHTTP(w, r)

			user := newEmptyUser()
			err := json.Unmarshal([]byte(w.Body.String()), user)
			if len(user.PassHash) != 0 || len(user.Email) != 0 || err != nil {
				t.Errorf("Expects no email or pass")
			}
			if w.Header().Get(headerContentType) != contentTypeJSON {
				t.Errorf("Expect %s but got %s", contentTypeJSON, w.Header().Get(headerContentType))
			}

		}
		handler.ServeHTTP(w, r)
		code := w.Code
		if code != c.expectedCode {
			t.Errorf("Expect %d but got %d", c.expectedCode, code)
		}
	}
}

func TestSpecificUserHandler(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Erorr connecting to db %v", err)
	}
	defer db.Close()
	userStore := users.NewMySQLStore(db)
	_, err = sessions.NewSessionID("somerandomkey")
	if err != nil {
		t.Fatalf("error creating new sessions %v", err)
	}
	ctx := &HandlerContext{
		Key:          "somerandomkey",
		SessionStore: sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), time.Hour),
		UserStore:    userStore,
	}
	cases := []struct {
		caseName      string
		requestType   string
		requestMethod string
		user          *users.User
		path          string
		update        *users.Updates
		expectedCode  int
		hasBearer     bool
		errorType     int
	}{
		{
			"Error: not authorized",
			contentTypeJSON,
			http.MethodPost,
			newEmptyUser(),
			"666666",
			&users.Updates{},
			http.StatusUnauthorized,
			false,
			0,
		},
		{
			"Error: method is forbidden",
			contentTypeJSON,
			http.MethodPost,
			newEmptyUser(),
			"666666",
			&users.Updates{},
			http.StatusMethodNotAllowed,
			true,
			1,
		},

		{
			"Error: user patch forbidden",
			contentTypeJSON,
			http.MethodPatch,
			&users.User{
				ID: 666666,
			},
			"1",
			&users.Updates{},
			http.StatusForbidden,
			true,
			2,
		},
		{
			"Error: request type incorrect",
			"text/plain",
			http.MethodPatch,
			&users.User{
				ID: 666666,
			},
			"1",
			&users.Updates{},
			http.StatusUnsupportedMediaType,
			true,
			2,
		},
		{
			"Error: user not found",
			contentTypeJSON,
			http.MethodGet,
			&users.User{
				ID: 666666,
			},
			"666666",
			&users.Updates{},
			http.StatusNotFound,
			true,
			1,
		},
		{
			"Error: update internal error",
			contentTypeJSON,
			http.MethodPatch,
			&users.User{
				ID: 666,
			},
			"666",
			&users.Updates{
				FirstName: "f",
			},
			http.StatusInternalServerError,
			true,
			2,
		},
		{
			"Error: update bad request",
			contentTypeJSON,
			http.MethodPatch,
			&users.User{
				ID: 1,
			},
			"1",
			&users.Updates{},
			http.StatusBadRequest,
			true,
			2,
		},
		{
			"Sucess: update user by id",
			contentTypeJSON,
			http.MethodPatch,
			newSpecificUser(666),
			"666",
			&users.Updates{
				FirstName: "f",
				LastName:  "l",
			},
			http.StatusOK,
			true,
			3,
		}, {
			"Sucess: get user with me url",
			contentTypeJSON,
			http.MethodGet,
			newSpecificUser(666),
			"me",
			&users.Updates{},
			http.StatusOK,
			true,
			3,
		},
		{
			"Sucess: get user",
			contentTypeJSON,
			http.MethodGet,
			newSpecificUser(666),
			"666",
			&users.Updates{},
			http.StatusOK,
			true,
			3,
		},
		{
			"Sucess: update user by me",
			contentTypeJSON,
			http.MethodPatch,
			newSpecificUser(666),
			"me",
			&users.Updates{
				FirstName: "f",
				LastName:  "l",
			},
			http.StatusOK,
			true,
			3,
		},
	}
	for _, c := range cases {
		body := []byte("{json:json}")
		if 0 != len(c.update.LastName) || 0 != len(c.update.FirstName) {
			body, err = json.Marshal(c.update)
			if err != nil {
				t.Errorf("json encode went wrong")
			}
		}
		handler := http.HandlerFunc(ctx.SpecificUserHandler)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(c.requestMethod, "/v1/users/"+c.path, strings.NewReader(string(body)))
		if err != nil {
			t.Errorf("Error %v", err)
		}
		r.Header.Set(headerContentType, c.requestType)
		sess := &SessionState{
			Time: time.Now(),
			User: c.user,
		}
		sid, err := sessions.BeginSession(ctx.Key, ctx.SessionStore, sess, w)
		if err != nil {
			t.Errorf("sid creation went wrong %v", err.Error())
		}
		r.Header.Set(authorization, "Bearer "+string(sid))
		if c.errorType == 3 {
			mock.ExpectQuery("select id, email, username, password_hash, first_name, last_name, photo_url from users where id=?").WithArgs(c.user.ID).WillReturnRows(mock.NewRows([]string{"ID", "Email", "PassHash", "UserName", "FirstName", "LastName", "PhotoURL"}).AddRow(c.user.ID, c.user.Email, c.user.PassHash, c.user.UserName, c.user.FirstName, c.user.LastName, c.user.PhotoURL))
			if http.MethodPatch == c.requestMethod {
				mock.ExpectExec("update users set first_name=?, last_name=? where id=?").WithArgs(c.update.FirstName, c.update.LastName, c.user.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			handler.ServeHTTP(w, r)
			if contentTypeJSON != w.Header().Get(headerContentType) {
				t.Errorf("Expect %v but got %v", contentTypeJSON, w.Header().Get(headerContentType))
			}
		}
		if c.errorType == 1 || c.errorType == 2 || c.errorType == 3 {
			handler.ServeHTTP(w, r)
			code := w.Code
			if code != c.expectedCode {
				t.Errorf("[Case %s] Expect %v but got %v; error %v", c.caseName, c.expectedCode, code, w.Body.String())
			}
		}
		if c.errorType == 2 || c.errorType == 3 {
			sess.User = c.user
			err = ctx.SessionStore.Save(sid, sess)
			if err != nil {
				t.Errorf("session save went wrong")
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations %s", err)
		}
	}
}

func TestSessionsHandler(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Erorr connecting to db %v", err)
	}
	defer db.Close()
	_, err = sessions.NewSessionID("somerandomkey")
	if err != nil {
		t.Fatalf("error creating new sessions %v", err)
	}
	ctx := &HandlerContext{
		Key:          "somerandomkey",
		SessionStore: sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), time.Hour),
		UserStore:    users.NewMySQLStore(db),
	}
	cases := []struct {
		caseName      string
		expectedCode  int
		requestType   string
		requestMethod string
		user          *users.User
		creds         *users.Credentials
		errorType     int
	}{
		{
			"Pass",
			http.StatusCreated,
			contentTypeJSON,
			http.MethodPost,
			newSpecificUser(666),
			&users.Credentials{
				Email:    "testmail@testmail.com",
				Password: "Password@123",
			},
			2,
		},
		{
			"Error: password authentication failed",
			http.StatusUnauthorized,
			contentTypeJSON,
			http.MethodPost,
			newSpecificUser(1),
			&users.Credentials{
				Email:    "testmail@testmail.com",
				Password: "Password@123",
			},
			1,
		},
		{
			"Error: requeest format incorrect",
			http.StatusUnsupportedMediaType,
			"text/plain",
			http.MethodPost,
			newEmptyUser(),
			&users.Credentials{},
			0,
		},
		{
			"Error: unsupported request method",
			http.StatusMethodNotAllowed,
			contentTypeJSON,
			http.MethodGet,
			newEmptyUser(),
			&users.Credentials{},
			0,
		},
		{
			"Error: user is not found",
			http.StatusUnauthorized,
			contentTypeJSON,
			http.MethodPost,
			newEmptyUser(),
			&users.Credentials{
				Email: "testmail@testmail.com",
			},
			0,
		},
		{
			"Error: credential is empty",
			http.StatusBadRequest,
			contentTypeJSON,
			http.MethodPost,
			newEmptyUser(),
			&users.Credentials{},
			0,
		},
	}
	for _, c := range cases {
		handler := http.HandlerFunc(ctx.SessionsHandler)
		cr := []byte("{json:json}")
		if 0 != len(c.creds.Password) || 0 != len(c.creds.Email) {
			cr, err = json.Marshal(c.creds)
			if err != nil {
				t.Errorf("json encode went wrong")
			}
		}
		w := httptest.NewRecorder()
		r, err := http.NewRequest(c.requestMethod, "/v1/sessions", strings.NewReader(string(cr)))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		r.Header.Set(headerContentType, c.requestType)
		if c.errorType == 3 || c.errorType == 2 {
			c.user.PassHash, err = bcrypt.GenerateFromPassword(c.user.PassHash, 13)
			if err != nil {
				t.Errorf("Hash creation went wrong")
			}
			mock.ExpectQuery("select id, email, username, password_hash, first_name, last_name, photo_url from users where email=?").WithArgs(c.user.Email).WillReturnRows(mock.NewRows([]string{"ID", "Email", "UserName", "PassHash", "FirstName", "LastName", "PhotoURL"}).AddRow(c.user.ID, c.user.Email, c.user.UserName, c.user.PassHash, c.user.FirstName, c.user.LastName, c.user.PhotoURL))

			handler.ServeHTTP(w, r)
			if w.Header().Get(headerContentType) != contentTypeJSON {
				t.Errorf("Expected a content type of [%s] got [%s]", contentTypeJSON, w.Header().Get(headerContentType))
			}
			user := newEmptyUser()
			err := json.Unmarshal([]byte(w.Body.String()), user)
			if err != nil || len(user.Email) != 0 || len(user.PassHash) != 0 {
				t.Errorf("Expect body to be empety")
			}
		}
		if c.errorType == 1 {
			mock.ExpectQuery("select id, email, username, password_hash, first_name, last_name, photo_url from users where email=?").WithArgs(c.user.Email).WillReturnRows(mock.NewRows([]string{"ID", "Email", "UserName", "PassHash", "FirstName", "LastName", "PhotoURL"}).AddRow(c.user.ID, c.user.Email, c.user.UserName, c.user.PassHash, c.user.FirstName, c.user.LastName, c.user.PhotoURL))
		}
		handler.ServeHTTP(w, r)
		code := w.Code
		if code != c.expectedCode {
			t.Errorf("[Case %s] Expect %d but got %d; message: %s", c.caseName, c.expectedCode, code, w.Body.String())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectation: %s", err)
		}
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	ctx := &HandlerContext{
		Key:          "somerandomkey",
		SessionStore: sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), time.Hour),
		UserStore:    &users.FakeMySQLStore{},
	}
	_, err := sessions.NewSessionID("somerandomkey")
	if err != nil {
		t.Errorf("problem creating newe session %v", err)
	}
	cases := []struct {
		caseName      string
		url           string
		expectedCode  int
		requestMethod string
		hasSession    bool
		errorType     int
	}{
		{
			"Pass",
			"mine",
			http.StatusMethodNotAllowed,
			http.MethodDelete,
			false,
			1,
		},
		{
			"Error: method not allowed",
			"mine",
			http.StatusMethodNotAllowed,
			http.MethodPatch,
			false,
			0,
		},
		{
			"Error: seesion can't be deleted",
			"mine",
			http.StatusInternalServerError,
			http.MethodDelete,
			true,
			0,
		},
		{
			"Error: incorrect url name",
			"dfasdfafsd",
			http.StatusForbidden,
			http.MethodDelete,
			false,
			0,
		},
	}
	for _, c := range cases {
		handler := http.HandlerFunc(ctx.SpecificSessionHandler)
		w := httptest.NewRecorder()
		sess := &SessionState{}
		sid, err := sessions.BeginSession(ctx.Key, ctx.SessionStore, sess, w)
		if err != nil {
			t.Fatal("Session begin error")
		}
		r, err := http.NewRequest(c.requestMethod, "/v1/sessions/"+c.url, nil)
		if err != nil {
			t.Errorf("Error %v", err)
		}
		r.Header.Set(authorization, "Bearer "+string(sid))
		if 1 == c.errorType {
			handler.ServeHTTP(w, r)
			expMsg := "signed out"
			actualMsg := w.Body.String()
			if expMsg != actualMsg {
				t.Errorf("%s: Expected a response of [%s] got [%s].", c.caseName, expMsg, actualMsg)
			}
		}
		if 0 == c.errorType {
			if c.hasSession {
				sid, err = sessions.NewSessionID("otherkey")
				if err != nil {
					t.Errorf("Cannot create new session")
				}
			}
			r.Header.Set(authorization, "Bearer "+string(sid))
			handler.ServeHTTP(w, r)
			code := w.Code
			if code != c.expectedCode {
				t.Errorf("[Case %s] Expect %d but got %d; Problem %s", c.caseName, c.expectedCode, code, w.Body.String())
			}
		}

	}
}
