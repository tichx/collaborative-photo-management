package users

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

func createNewUser() *NewUser {
	return &NewUser{
		UserName:     "info441",
		Email:        "testmail@testmail.com",
		FirstName:    "First",
		LastName:     "Last",
		Password:     "Password@123",
		PasswordConf: "Password@123",
	}
}

func TestValidate(t *testing.T) {
	user := createNewUser()
	err := user.Validate()
	if err != nil {
		t.Errorf("Validation failed. Got: %v", err)
	}
	user.UserName = ""
	err = user.Validate()
	if err.Error() != "UserName must be non-zero length and may not contain spaces" {
		t.Errorf("Expecting an error username format, instead got: %v", err)
	}
	user.UserName = "info441"

	user.UserName = "info 441"
	err = user.Validate()
	if err.Error() != "UserName must be non-zero length and may not contain spaces" {
		t.Errorf("Expecting an error username format, instead got: %v", err)
	}
	user.UserName = "info441"

	user.Email = "dasfsdf"
	err = user.Validate()
	if err.Error() != "Invalid email address" {
		t.Errorf("Expecting an error email, instead got: %v", err)
	}
	user.Email = "testmail@testmail.com"

	user.Password = "o"
	err = user.Validate()
	if err.Error() != "Password must be at least 6 characters" {
		t.Errorf("Expecting an error password format, instead got: %v", err)
	}
	user.Password = "Password@123"

	user.PasswordConf = "password@123"
	err = user.Validate()
	if err.Error() != "Password and Password Confirmation must match" {
		t.Errorf("Expecting an error email match, instead got: %v", err)
	}
	user.PasswordConf = "Password@123"

}

func TestToUser(t *testing.T) {
	user := createNewUser()
	userConverted, err := user.ToUser()
	if err != nil {
		t.Errorf("Invalid Conversion: %v", err.Error())

	}
	if userConverted == nil {
		t.Errorf("Empty user")
	}
	if user.Email != userConverted.Email {
		t.Errorf("Email Doesn't match. Expected: %v; Got: %v", user.Email, userConverted.Email)
	}
	if len(userConverted.PassHash) == 0 {
		t.Errorf("Invalid PassHash")
	}
	err = bcrypt.CompareHashAndPassword(userConverted.PassHash, []byte(user.Password))
	if err != nil {
		t.Errorf("Invalid Passhash")

	}
	if len(userConverted.PhotoURL) == 0 || !strings.HasPrefix(userConverted.PhotoURL, gravatarBasePhotoURL) {
		t.Errorf("User.PhotoURL is zero length, should be gravatar profile image URL\n")
	}
	invalid := createNewUser()
	invalid.Email = "invalid@@@"
	_, err = invalid.ToUser()
	if err == nil {
		t.Errorf("Expected: invalid error; Got: none")
	}
	invalid.Email = " TestMail@TestMail.com"
	userConverted, err = invalid.ToUser()
	if err != nil {
		t.Errorf("Unexpected conversion error")

	}
	user = createNewUser()
	u, err := user.ToUser()
	if err != nil {
		t.Errorf("Unexpected conversion error")
	}
	if userConverted.PhotoURL != u.PhotoURL {
		t.Errorf("PhotoURL does not match. Expected: %v. Got: %v", u.PhotoURL, userConverted.PhotoURL)
	}
}

func TestFullName(t *testing.T) {
	newUser := createNewUser()
	user, err := newUser.ToUser()
	if err != nil {
		t.Errorf("Error converting new user")
	}

	user.FirstName = ""
	user.LastName = ""
	name := user.FullName()
	if name != "" {
		t.Errorf("Expected: . Got %v", name)
	}

	user.LastName = "Last"
	name = user.FullName()
	if name != "Last" {
		t.Errorf("Expected: Last. Got %v", name)

	}
	user.FirstName = "First"
	name = user.FullName()
	if name != "First Last" {
		t.Errorf("Expected: First. Got %v", name)
	}

	user.LastName = ""
	name = user.FullName()
	if name != "First" {
		t.Errorf("Expected: First Last. Got %v", name)
	}

}

func TestAuthenticate(t *testing.T) {
	user, err := createNewUser().ToUser()
	if err != nil {
		t.Errorf("Error converting new user")

	}
	err = user.Authenticate("")
	if err == nil {
		t.Errorf("Error in authentication: empty password should not be authenticated")
	}
	err = user.SetPassword("newpassword")
	if err != nil {
		t.Errorf("Error: %v", err)

	}
	err = user.Authenticate("wrongpassword")
	if err == nil {
		t.Errorf("Error in authentication: wrong password should not be authenticated")
	}
	err = user.Authenticate("newpassword")
	if err != nil {
		t.Errorf("Error in autentication: %v", err)

	}

}

func TestApplyUpdates(t *testing.T) {
	user, err := createNewUser().ToUser()
	if err != nil {
		t.Errorf("Error converting new user")

	}

	updates := &Updates{
		FirstName: "",
		LastName:  "",
	}
	err = user.ApplyUpdates(updates)
	if err == nil {
		t.Errorf("Expected an error when updates are empty.")
	}
	if user.FirstName != "First" && user.LastName != "Last" {
		t.Errorf("Expected: %s %s. Got: %s, %s", "First", "Last", user.FirstName, user.LastName)

	}

	updates.FirstName = "first"
	err = user.ApplyUpdates(updates)
	if err != nil {
		t.Errorf("Error in apply updates: %v", err)

	}
	if user.FirstName != "first" && user.LastName != "Last" {
		t.Errorf("Expected: %s %s. Got: %s %s", "first", "Last", user.FirstName, user.LastName)
	}

	updates.LastName = "last"
	err = user.ApplyUpdates(updates)
	if err != nil {
		t.Errorf("Error in apply updates: %v", err)

	}
	if user.FirstName != "first" && user.LastName != "last" {
		t.Errorf("Expected: %s %s. Got: %s %s", "first", "last", user.FirstName, user.LastName)
	}

	updates.LastName = ""
	err = user.ApplyUpdates(updates)
	if err != nil {
		t.Errorf("Error in apply updates: %v", err)

	}
	if user.FirstName != "first" && user.LastName != "" {
		t.Errorf("Names not updated properly. Expected: %s %s. Got: %s, %s", "first", "", user.FirstName, user.LastName)

	}

}
