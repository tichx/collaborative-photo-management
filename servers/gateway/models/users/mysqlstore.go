package users

import (
	"database/sql"
	"fmt"
)

//MySQLStorego is a connection to the database
type MySQLStorego struct {
	db *sql.DB
}

//SignInLog is a struct for user sign in logs
type SignInLog struct {
	ID       string
	DateTime string
	IPAddr   string
}

//NewMySQLStorego Constructs new MySQLStorego
func NewMySQLStorego(db *sql.DB) *MySQLStorego {
	return &MySQLStorego{
		db: db,
	}
}

const getByIDPrepared = "SELECT id, email, first_name, last_name, user_name, pass_hash, photo_url FROM User WHERE id=?"
const getByEmailPrepared = "SELECT id, email, first_name, last_name, user_name, pass_hash, photo_url FROM User WHERE email=?"
const getByUserNamePrepared = "SELECT id, email, first_name, last_name, user_name, pass_hash, photo_url FROM User WHERE user_name=?"
const insertPrepared = "INSERT INTO User(email, first_name, last_name, user_name, pass_hash, photo_url) values (?,?,?,?,?,?)"
const logSigninPrepared = "insert into UserSignIns(id, date_time, ip_addr) values (?,?,?)"
const deletePrepared = "DELETE FROM User WHERE id = ?"

// const updatePrepared = "UPDATE User SET first_name = ?, last_name = ? WHERE id = ?"
const updatePrepared = "UPDATE User SET first_name = ?, last_name = ? WHERE id = ?"

//GetByID returns the User with the given ID
func (ms *MySQLStorego) GetByID(id int64) (*User, error) {
	user := User{}
	userRow, err := ms.db.Query(getByIDPrepared, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	defer userRow.Close()
	for userRow.Next() {
		err = userRow.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("err scanning user row: %v", err)
		}
	}
	err = userRow.Err()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//GetByEmail returns the User with the given email
func (ms *MySQLStorego) GetByEmail(email string) (*User, error) {
	user := User{}
	userRow, err := ms.db.Query(getByEmailPrepared, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	defer userRow.Close()
	for userRow.Next() {
		err = userRow.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("err scanning user row: %v", err)
		}
	}
	err = userRow.Err()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//GetByUserName returns the User with the given Username
func (ms *MySQLStorego) GetByUserName(username string) (*User, error) {
	user := User{}
	userRow, err := ms.db.Query(getByUserNamePrepared, username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	defer userRow.Close()
	for userRow.Next() {
		err = userRow.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("err scanning user row: %v", err)
		}
	}
	err = userRow.Err()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ms *MySQLStorego) Insert(user *User) (*User, error) {
	result, err := ms.db.Exec(insertPrepared, &user.Email, &user.FirstName, &user.LastName, &user.UserName, &user.PassHash, &user.PhotoURL)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error user new id: %v", err)
	}

	user.ID = insertID
	return user, nil
}

//LogSignin inserts the successful user sign in into the database, and returns
//the newly-inserted log
func (ms *MySQLStorego) LogSignin(logSignin *SignInLog) (*SignInLog, error) {

	res, err := ms.db.Exec(logSigninPrepared, logSignin.ID, logSignin.DateTime, logSignin.IPAddr)
	if err != nil {
		return nil, fmt.Errorf("error inserting log: %v", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting inserted log: %v", err)
	}
	return logSignin, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *MySQLStorego) Update(id int64, updates *Updates) (*User, error) {
	updateSQL := updatePrepared
	_, err := ms.db.Exec(updateSQL, updates.FirstName, updates.LastName, id)
	if err != nil {
		return nil, fmt.Errorf("error update user: %v", err)
	}
	updateUser, err := ms.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting updating user by id: %v", err)
	}
	return updateUser, nil
}

//Delete deletes the user with the given ID
func (ms *MySQLStorego) Delete(id int64) error {
	_, err := ms.db.Exec(deletePrepared, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}
