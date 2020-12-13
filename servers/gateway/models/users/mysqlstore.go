package users

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//MySQLStore will be used in the func
type MySQLStore struct {
	Client *sql.DB
}

//NewMySQLStore returns a new
func NewMySQLStore(db *sql.DB) *MySQLStore {
	return &MySQLStore{db}
}

//GetByID returns the User with the given ID
func (db *MySQLStore) GetByID(id int64) (*User, error) {
	user := &User{}
	qs := "select id, email, username, password_hash, first_name, last_name, photo_url from users where id=?"
	row := db.Client.QueryRow(qs, id)
	if err := row.Scan(&user.ID, &user.Email, &user.UserName, &user.PassHash,
		&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//GetByEmail returns the User with the given email
func (db *MySQLStore) GetByEmail(email string) (*User, error) {
	sq := "select id, email, username, password_hash, first_name, last_name, photo_url from users where email=?"
	row := db.Client.QueryRow(sq, email)
	user := &User{}
	if err := row.Scan(&user.ID, &user.Email, &user.UserName, &user.PassHash,
		&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//GetByUserName returns the User with the given Username
func (db *MySQLStore) GetByUserName(username string) (*User, error) {
	sq := "select id, email, username, password_hash, first_name, last_name, photo_url from users where username=?"
	row := db.Client.QueryRow(sq, username)
	user := &User{}
	if err := row.Scan(&user.ID, &user.Email, &user.UserName, &user.PassHash,
		&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (db *MySQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, username, password_hash, first_name, last_name, photo_url) values (?,?,?,?,?,?)"
	msg, err := db.Client.Exec(insq, &user.Email, &user.UserName, &user.PassHash,
		&user.FirstName, &user.LastName, &user.PhotoURL)
	if err != nil {
		return nil, err
	}
	id, err := msg.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (db *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	user, err := db.GetByID(id)
	if err != nil {
		return nil, err
	}
	err = user.ApplyUpdates(updates)
	if err != nil {
		return nil, err
	}
	upd := "update users set first_name=?, last_name=? where id=?"
	_, err = db.Client.Exec(upd, user.FirstName, user.LastName, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil

}

//InsertSignIn logs a new successful sign in
func (db *MySQLStore) InsertSignIn(userID int64, ip string) (int64, error) {
	insq := "insert into users_signin(id, date_time, ip_addr) values (?,now(),?)"
	res, err := db.Client.Exec(insq, userID, ip)
	if err != nil {
		return int64(0), err
	}

	//get generated ID from insert
	id, err := res.LastInsertId()
	if err != nil {
		return int64(0), err
	}
	return id, nil
}

//Delete the user with the given ID
func (db *MySQLStore) Delete(id int64) error {
	dq := "delete from users where id=?"
	_, err := db.Client.Exec(dq, id)
	if err != nil {
		return err
	}
	return nil
}
