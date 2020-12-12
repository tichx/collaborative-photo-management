package users

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetByID(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User with wrong id",
			&User{
				ID:        1,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			999,
			false,
		},
		{
			"GetByID Sucessful",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			10000,
			false,
		},
		{
			"Non-existed user",
			&User{},
			2,
			true,
		},
	}
	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := MySQLStore{db}

		row := mock.NewRows([]string{
			"ID",
			"Email",
			"UserName",
			"PassHash",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.UserName,
			c.expectedUser.PassHash,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id, email, username, password_hash, first_name, last_name, photo_url from users where id=?"
		if !c.expectError {
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}

		} else {
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}

		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByEmail(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"GetByEmail Success",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			"test@test.com",
			false,
		},
		{
			"GetByEmail Fail Non-existed user",
			&User{},
			"test123@test.com",
			true,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStore{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"UserName",
			"PassHash",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.UserName,
			c.expectedUser.PassHash,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id, email, username, password_hash, first_name, last_name, photo_url from users where email=?"
		if !c.expectError {
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}

		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestGetByUserName(t *testing.T) {
	cases := []struct {
		name          string
		expectedUser  *User
		userNameToGet string
		expectError   bool
	}{
		{
			"User Found",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			"username",
			false,
		},
		{
			"Username as Number Found",
			&User{
				ID:        100001,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username1",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			"username1",
			false,
		},
		{
			"Empty User Not Found",
			&User{},
			"username1",
			true,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStore{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"UserName",
			"PassHash",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.UserName,
			c.expectedUser.PassHash,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id, email, username, password_hash, first_name, last_name, photo_url from users where username=?"
		if !c.expectError {
			mock.ExpectQuery(query).WithArgs(c.userNameToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByUserName(c.userNameToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.userNameToGet).WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByUserName(c.userNameToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}

		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestInsert(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		expectError  bool
	}{
		{
			"Insert Fail: require a password",
			&User{
				Email:     "test3@test.com",
				UserName:  "username13",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			true,
		},
		{
			"Insert Fail: require a username",
			&User{
				Email:     "test2@test.com",
				PassHash:  []byte("passhash123"),
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			true,
		},
		{
			"Insert Fail: require an email",
			&User{
				PassHash:  []byte("passhash123"),
				UserName:  "username12",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			true,
		},
		{
			"Insert Succuss",
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username1",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &MySQLStore{db}
		exec := "insert into users(email, username, password_hash, first_name, last_name, photo_url) values (?,?,?,?,?,?)"
		if !c.expectError {

			mock.ExpectExec(exec).WithArgs(c.expectedUser.Email, c.expectedUser.UserName, c.expectedUser.PassHash, c.expectedUser.FirstName, c.expectedUser.LastName, c.expectedUser.PhotoURL).WillReturnResult(sqlmock.NewResult(1, 1))
			user, err := mainSQLStore.Insert(c.expectedUser)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		} else {
			mock.ExpectExec(exec).WithArgs(c.expectedUser.Email, c.expectedUser.UserName, c.expectedUser.PassHash, c.expectedUser.FirstName, c.expectedUser.LastName, c.expectedUser.PhotoURL).WillReturnError(ErrUserNotFound)
			user, err := mainSQLStore.Insert(c.expectedUser)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name         string
		originalUser *User
		updates      *Updates
		expectedUser *User
		expectError  bool
	}{
		{
			"Parse empty updates struct",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			&Updates{},
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			true,
		},
		{
			"Update first",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			&Updates{
				FirstName: "F",
			},
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "F",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			false,
		},
		{
			"Update last",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			&Updates{
				LastName: "L",
			},
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "L",
				PhotoURL:  "photourl",
			},
			false,
		},
		{
			"Update first and last",
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			&Updates{
				FirstName: "First",
				LastName:  "Last",
			},
			&User{
				ID:        10000,
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "First",
				LastName:  "Last",
				PhotoURL:  "photourl",
			},
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &MySQLStore{db}
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"UserName",
			"PassHash",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.UserName,
			c.expectedUser.PassHash,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		exec := "update users set first_name=?, last_name=? where id=?"
		query := "select id, email, username, password_hash, first_name, last_name, photo_url from users where id=?"
		if !c.expectError {

			mock.ExpectQuery(query).WithArgs(c.originalUser.ID).WillReturnRows(row)
			if len(c.updates.FirstName) == 0 {
				mock.ExpectExec(exec).WithArgs(c.originalUser.FirstName, c.updates.LastName, c.originalUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			} else if len(c.updates.LastName) == 0 {
				mock.ExpectExec(exec).WithArgs(c.updates.FirstName, c.originalUser.LastName, c.originalUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(exec).WithArgs(c.updates.FirstName, c.updates.LastName, c.originalUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			}
			user, err := mainSQLStore.Update(c.originalUser.ID, c.updates)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		} else {
			mock.ExpectQuery(query).WithArgs(c.originalUser.ID).WillReturnRows(row)
			user, err := mainSQLStore.Update(c.originalUser.ID, c.updates)
			if user != nil {
				t.Errorf("Expected error [%v] but got [%v] instead", errors.New("Update went wrong"), err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name         string
		expectedUser *User
		idToDelete   int64
		expectError  bool
	}{
		{
			"Delete failed by not supplying id",
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			-1,
			true,
		},
		{
			"Delete failed by invalid id",
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			2913801231038,
			true,
		},
		{
			"Delete success",
			&User{
				Email:     "test@test.com",
				PassHash:  []byte("passhash123"),
				UserName:  "username",
				FirstName: "firstname",
				LastName:  "lastname",
				PhotoURL:  "photourl",
			},
			1,
			false,
		},
	}

	for _, c := range cases {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()
		mainSQLStore := &MySQLStore{db}
		dq := "delete from users where id=?"
		sq := "select id, email, username, password_hash, first_name, last_name, photo_url from users where id=?"
		if !c.expectError {
			mock.ExpectExec(dq).WithArgs(c.idToDelete).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery(sq).WithArgs(c.idToDelete).WillReturnError(ErrUserNotFound)
			err := mainSQLStore.Delete(c.idToDelete)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			user, err := mainSQLStore.GetByID(c.idToDelete)
			if user != nil || err == nil {
				t.Errorf("The id assoicated with the user is not deleted")
			}
		} else {
			mock.ExpectExec(dq).WithArgs(c.idToDelete).WillReturnError(ErrUserNotFound)
			err := mainSQLStore.Delete(c.idToDelete)
			if err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}
	}
}
