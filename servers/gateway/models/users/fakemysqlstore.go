package users

// FakeMySQLStore stores fake users
type FakeMySQLStore struct {
}

//NewFakeConn returns a new connection object
func NewFakeConn() *FakeMySQLStore {
	return &FakeMySQLStore{}
}

//GetByEmail gives user by email
func (store *FakeMySQLStore) GetByEmail(email string) (*User, error) {
	return nil, nil

}

//GetByID gives the user by id
func (store *FakeMySQLStore) GetByID(id int64) (*User, error) {
	return nil, nil
}

//GetByUserName gives users by username
func (store *FakeMySQLStore) GetByUserName(username string) (*User, error) {
	return nil, nil

}

//Delete removes users by id
func (store *FakeMySQLStore) Delete(id int64) error {
	return nil
}

//Update updates fake users
func (store *FakeMySQLStore) Update(id int64, updates *Updates) (*User, error) {
	return nil, nil
}

//Insert sets database with new user
func (store *FakeMySQLStore) Insert(user *User) (*User, error) {
	if user.UserName == "username1" {
		return nil, ErrUserNotFound
	} else if user.ID == 0 {
		user.ID = 1
		return user, nil
	}
	return nil, ErrUserNotFound
}
