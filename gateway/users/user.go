package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	//validate the new user according to these rules:

	//- Email field must be a valid email address, also doesn't contain space. Basic check
	_, err := mail.ParseAddress(nu.Email)
	nu.Email = strings.ToLower(strings.Trim(nu.Email, " "))
	if err != nil || strings.Contains(nu.Email, " ") {
		return fmt.Errorf("error: not a valid email address")
	}
	// //A cool regex email check that cannot be used for this assignment due to its design
	// //but worth keeping here as an option
	// var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(nu.Email) > 255 {
		return fmt.Errorf("error: email address too long. Stop Attacking :X")
	}

	//- Password must be at least 6 characters
	//- Password and PasswordConf must match
	if len(nu.Password) < 6 {
		return fmt.Errorf("error: password too short")
	} else if nu.Password != nu.PasswordConf {
		return fmt.Errorf("error: password doesn't match")
	}
	//- UserName must be non-zero length and may not contain spaces
	if len(nu.UserName) == 0 || strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("error: invalid username")
	}
	//use fmt.Errorf() to generate appropriate error messages if
	//the new user doesn't pass one of the validation rules

	/*
		//A possible add-on function to use regex match convention
		fnMatched, err := regexp.MatchString(`^[\w'\-,.][^0-9_!¡?÷?¿/\\+=@#$%ˆ&*(){}|~<>;:[\]]{2,}$`, nu.FirstName)
		if !fnMatched || err != nil {
			return fmt.Errorf("error: first name not in convention")
		}
	*/
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	//call Validate() to validate the NewUser and
	//return any validation errors that may occur.

	//if valid, create a new *User and set the fields
	//based on the field values in `nu`.
	//Leave the ID field as the zero-value; your Store
	//implementation will set that field to the DBMS-assigned
	//primary key value.
	//Set the PhotoURL field to the Gravatar PhotoURL
	//for the user's email address.
	//see https://en.gravatar.com/site/implement/hash/
	//and https://en.gravatar.com/site/implement/images/

	//TODO: also call .SetPassword() to set the PassHash
	//field of the User to a hash of the NewUser.Password
	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	h := md5.New()
	h.Write([]byte(strings.TrimSpace(strings.ToLower(nu.Email))))
	emailhash := h.Sum(nil)
	imageURL := gravatarBasePhotoURL + hex.EncodeToString(emailhash)
	user := &User{
		ID:        0,
		Email:     nu.Email,
		PassHash:  nil,
		UserName:  nu.UserName,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		PhotoURL:  imageURL,
	}
	user.SetPassword(nu.Password)

	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	fn := u.FirstName
	ln := u.LastName

	//if first name is empty
	if len(fn) == 0 {
		//if last name is also empty
		if len(ln) == 0 {
			return ""
		}
		//only first name is empty
		return ln
		// if first name not empty and last name empty
	} else if len(ln) == 0 {
		return fn
	}

	return (fn + " " + ln)
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password string cannot be less than 6 characters")
	}
	//use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	u.PassHash = hashedPW
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	//set the fields of `u` to the values of the related
	//field in the `updates` struct

	// // code to check if first or last name is empty. Sounds important but cannot be used for this assignment
	// if len(u.FirstName) == 0 || len(u.LastName) == 0 {
	// 	return fmt.Errorf("error: name cannot be empty")
	// }

	if updates == nil {
		return fmt.Errorf("updates struct is nil")
	}

	if updates.FirstName == "" || updates.LastName == "" {
		return fmt.Errorf("one or more updates fields is nil")
	}

	if u.FirstName == updates.FirstName && u.LastName == updates.LastName {
		return fmt.Errorf("error: first name and last name still the same")
	}

	u.FirstName = updates.FirstName
	u.LastName = updates.LastName

	return nil
}
