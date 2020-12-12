package users

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the MySQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStorego{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := getByIDPrepared

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).
				WithArgs(c.idToGet).
				WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByEmail is a test function for the MySQLStore's GetByEmail
func TestGetByEmail(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User Not Found",
			&User{},
			"test2@test.com",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStorego{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := getByEmailPrepared

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).
				WithArgs(c.emailToGet).
				WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByUserName is a test function for the MySQLStore's GetByUserName
func TestGetByUserName(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name          string
		expectedUser  *User
		usernameToGet string
		expectError   bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User Not Found",
			&User{},
			"username",
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStorego{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := getByUserNamePrepared

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).
				WithArgs(c.usernameToGet).
				WillReturnError(ErrUserNotFound)

			user, err := mainSQLStore.GetByUserName(c.usernameToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.usernameToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByUserName(c.usernameToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestInsert is a test function for the MySQLStore's Insert
func TestInsert(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *NewUser
		expectError  bool
	}{
		{
			"User Inserted",
			&NewUser{
				Email:        "test@test.com",
				Password:     "passhash123",
				PasswordConf: "passhash123",
				UserName:     "username",
				FirstName:    "firstname",
				LastName:     "lastname",
			},
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStorego{db}

		user, err := c.expectedUser.ToUser()
		if err != nil {
			t.Fatalf("unexpected error constructing user: %v", err)
		}

		query := insertPrepared
		var insertedID int64 = 1
		// Set up expected query that will expect an error
		mock.ExpectQuery(query).
			WithArgs(user.Email,
				user.FirstName,
				user.LastName,
				user.UserName,
				user.PassHash,
				user.PhotoURL).
			WillReturnError(ErrUserNotFound)

		res, err := mainSQLStore.Insert(user)
		if res != nil || err == nil {
			t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
		}

		// Set up an expected query with the expected result from the mock DB
		mock.ExpectExec(query).WithArgs(user.Email,
			user.FirstName,
			user.LastName,
			user.UserName,
			user.PassHash,
			user.PhotoURL).
			WillReturnResult(sqlmock.NewResult(insertedID, 1))

		// res, err = mainSQLStore.Insert(user)
		// if err != nil {
		// 	t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
		// }

	}
}

// TestDelete is a test function for the MySQLStore's Delete
func TestDelete(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name        string
		deleteUser  *User
		idToDelete  int64
		expectError bool
	}{
		{
			"User Deleted",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"Wrong User",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			2,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStorego{db}

		query := deletePrepared

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectExec(query).
				WithArgs(c.idToDelete).
				WillReturnError(ErrUserNotFound)

			// Test Delete()
			err := mainSQLStore.Delete(c.idToDelete)
			if err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectExec(query).WithArgs(c.idToDelete).WillReturnResult(sqlmock.NewResult(1, 1))

			// Test Delete()
			err := mainSQLStore.Delete(c.idToDelete)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestUpdateis a test function for the MySQLStore's Update
func TestUpdate(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		update       *Updates
		updatedUser  *User
		idToUpdate   int64
		expectError  bool
		blankUpdate  bool
	}{
		{
			"User Updated",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{
				FirstName: "upfirstname",
				LastName:  "uplastname",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"upfirstname",
				"uplastname",
				"photourl",
			},
			1,
			false,
			false,
		},
		{
			"User Updated Lastname",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{
				LastName: "uplastname",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"uplastname",
				"photourl",
			},
			1,
			false,
			false,
		},
		{
			"wrong id so not Updated",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{
				FirstName: "upfirstname",
				LastName:  "uplastname"},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			666,
			true,
			false,
		},
		{
			"Empty update so not Updated",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&Updates{},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			true,
			true,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := &MySQLStorego{db}

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.updatedUser.ID,
			c.updatedUser.Email,
			c.updatedUser.PassHash,
			c.updatedUser.UserName,
			c.updatedUser.FirstName,
			c.updatedUser.LastName,
			c.updatedUser.PhotoURL,
		)

		query := regexp.QuoteMeta(updatePrepared)
		queryID := getByIDPrepared

		if c.expectError {
			if c.blankUpdate {
				// Set up expected query that will expect an error
				// mock.ExpectQuery(query).
				//  WithArgs(c.update.FirstName, c.update.LastName, c.idToUpdate).
				//  WillReturnRows(row)
				mock.ExpectQuery(queryID).
					WithArgs(c.idToUpdate).
					WillReturnRows(row)
				// Test Update()
				user, err := mainSQLStore.Update(c.idToUpdate, c.update)
				if user != nil || err == nil {
					t.Errorf("Expected error [%v] but got [%v] instead", errors.New("updates can't be blank"), err)
				}
			} else {
				// Set up expected query that will expect an error
				mock.ExpectQuery(queryID).
					WithArgs(c.idToUpdate).
					WillReturnError(ErrUserNotFound)

				// Test Update()
				user, err := mainSQLStore.Update(c.idToUpdate, c.update)
				if user != nil || err == nil {
					t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
				}
			}

		} else {
			// queryID := getByIDPrepared
			// Set up an expected query with the expected result from the mock DB
			// mock.ExpectQuery(queryID).
			//  WithArgs(c.idToUpdate).
			//  WillReturnRows(row)

			if len(c.update.LastName) == 0 {
				mock.ExpectExec(query).WithArgs(c.update.FirstName, c.expectedUser.LastName, c.idToUpdate).WillReturnResult(sqlmock.NewResult(1, 1))
			} else if len(c.update.FirstName) == 0 {
				mock.ExpectExec(query).WithArgs(c.expectedUser.FirstName, c.update.LastName, c.idToUpdate).WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec(query).WithArgs(c.update.FirstName, c.update.LastName, c.idToUpdate).WillReturnResult(sqlmock.NewResult(1, 1))
			}

			// mock.ExpectExec(query).
			//  WithArgs(c.update.FirstName, c.update.LastName, c.idToUpdate).
			//  WillReturnResult(sqlmock.NewResult(1, 1))
			// Test UPdate()
			_, err := mainSQLStore.Update(c.idToUpdate, c.update)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s][%x][%x]: %v", c.name, c.expectedUser.ID, c.idToUpdate, err)
			}
			// if !reflect.DeepEqual(user, c.expectedUser) {
			//  t.Errorf("Error, invalid match in test [%s]", c.name)
			// }
		}

		// if err := mock.ExpectationsWereMet(); err != nil {
		//  t.Errorf("There were unfulfilled expectations: %s", err)
		// }

	}
}
