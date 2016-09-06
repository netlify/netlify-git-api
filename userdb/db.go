package userdb

import (
	"io/ioutil"
	"os"

	"github.com/pborman/uuid"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

// User is a user in the db
type User struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Email        string `yaml:"email"`
	PasswordHash string `yaml:"hash"`
}

// UserDB is the full set of users
type UserDB struct {
	dbPath string
	Users  []User `yaml:"users"`
}

// Read a userDB from a filepath
func Read(dbPath string) (*UserDB, error) {
	db := &UserDB{dbPath: dbPath}
	data, err := ioutil.ReadFile(dbPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	err = yaml.Unmarshal(data, db)
	return db, err
}

// Write the userDB
func (db *UserDB) Write() error {
	yml, err := yaml.Marshal(&db)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(db.dbPath, yml, 0644)
}

// LookupByEmail a user by email
func (db *UserDB) LookupByEmail(email string) *User {
	for _, user := range db.Users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

// Get a user by ID
func (db *UserDB) Get(id string) *User {
	for _, user := range db.Users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

// Add a new user to the db
func (db *UserDB) Add(email, name, pw string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := db.LookupByEmail(email)
	if user == nil {
		db.Users = append(db.Users, User{ID: uuid.New(), Name: name, Email: email, PasswordHash: string(hash)})
	} else {
		user.Name = name
		user.PasswordHash = string(hash)
	}

	return user, nil
}

// Delete a user from the db
func (db *UserDB) Delete(email string) {
	users := []User{}
	for _, user := range db.Users {
		if user.Email != email {
			users = append(users, user)
		}
	}
	db.Users = users
}

// Authenticate checks if a password is valid for this user
func (u *User) Authenticate(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pw))
	return err == nil
}
