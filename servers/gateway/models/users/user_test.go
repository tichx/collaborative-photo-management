package users

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// catches all possible validation errors, and returns
// no error when the new user is valid
func TestValidate(t *testing.T) {
	cases := []struct {
		nu          NewUser
		expectError bool
	}{
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			false,
		},
		{
			NewUser{
				Email:        "invalidemail",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "advalidlstnadadad",
			},
			true,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "mismatchpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			true,
		},

		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			true,
		},

		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "space usern",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			true,
		},

		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "spw",
				PasswordConf: "spw",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			true,
		},
	}
	for _, c := range cases {
		err := c.nu.Validate()
		if c.expectError && err == nil {
			t.Errorf("expected error")
		} else if !c.expectError && err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}

}

//calculates the PhotoURL field correctly, even when the
//email address has upper case letters or spaces (Links to
//an external site.), and sets the PassHash field to the password hash
func TestToUser(t *testing.T) {
	cases := []struct {
		nu          NewUser
		expectError bool
	}{
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			false,
		},
		{
			NewUser{
				Email:        " Valid@uw.com ",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			false,
		},
	}
	for _, c := range cases {
		u, err := c.nu.ToUser()
		if c.expectError && err == nil {
			t.Errorf("expected error")
		} else if !c.expectError && err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		hasher := md5.New()
		hasher.Write([]byte("valid@uw.com"))
		correctHashString := hex.EncodeToString(hasher.Sum(nil))
		right := "https://www.gravatar.com/avatar/" + correctHashString
		if u.PhotoURL != right {
			t.Errorf("incorrect photoURL")
		}
		if err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(c.nu.Password)); err != nil {
			t.Errorf("password hash not matched")
		}
	}

}

//verify that it returns the correct results given the
//various possible inputs
func TestFullName(t *testing.T) {
	cases := []struct {
		nu          NewUser
		correctName string
		expectError bool
	}{
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			"validfrstn validlstn",
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "",
				LastName:     "validlstn",
			},
			"validlstn",
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "",
			},
			"validfrstn",
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "",
				LastName:     "",
			},
			"",
			false,
		},
	}

	for _, c := range cases {
		u, err := c.nu.ToUser()
		if !c.expectError && err != nil {
			t.Errorf("Unexpect Error %v", err)
		}
		if u.FullName() != c.correctName {
			t.Errorf("fullname not matched")
		}
	}
}

//verify that authentication happens correctly
//for the various possible inputs
func TestAuthenticate(t *testing.T) {
	cases := []struct {
		nu          NewUser
		input       string
		expectError bool
	}{
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			"validpassword",
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			"wrongpassword",
			true,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			" ",
			true,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "Valid!@#$%^/\\pwd",
				PasswordConf: "Valid!@#$%^/\\pwd",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			"Valid!@#$%^/\\pwd",
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			"",
			true,
		},
	}

	for _, c := range cases {
		u, err := c.nu.ToUser()
		if err != nil {
			t.Errorf("Unexpect Error %v", err)
		}
		inputTest := u.Authenticate(c.input)
		if !c.expectError && inputTest != nil {
			t.Errorf("Unexpect Error %v", inputTest)
		}
		if c.expectError && inputTest == nil {
			t.Errorf("expected error")
		}
	}

}

//ensure the user's fields are updated properly given an Updates struct.
func TestApplydates(t *testing.T) {
	cases := []struct {
		nu          NewUser
		update      Updates
		expectError bool
	}{
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			Updates{
				FirstName: "",
				LastName:  "",
			},
			true,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			Updates{
				FirstName: "updatefrstn",
				LastName:  "",
			},
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			Updates{
				FirstName: "",
				LastName:  "updatelstn",
			},
			false,
		},
		{
			NewUser{
				Email:        "valid@uw.com",
				Password:     "validpassword",
				PasswordConf: "validpassword",
				UserName:     "validusrname",
				FirstName:    "validfrstn",
				LastName:     "validlstn",
			},
			Updates{
				FirstName: "updatefrstn",
				LastName:  "updatelstn",
			},
			false,
		},
	}

	for _, c := range cases {
		u, err := c.nu.ToUser()
		if err != nil {
			t.Errorf("unexpected Error %v", err)
		}
		inputTest := u.ApplyUpdates(&c.update)
		if !c.expectError && inputTest != nil {
			t.Errorf("unexpected Error %v", err)
		}
		if c.expectError && inputTest == nil {
			t.Errorf("expected error")
		}
		if inputTest == nil {
			if u.FirstName != c.update.FirstName || u.LastName != c.update.LastName {
				t.Errorf("did not update, original %s expected %s", u.FirstName, c.update.FirstName)
			}
		}

	}

}
