package boltdb

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// usersBucket stores the name of the boltdb bucket the users are stored in
var usersBucket = []byte("users")

// UserService contains the implemented functionality for users
type UserService struct {
	db *bolt.DB
}

// GetAll returns all users
func (us UserService) GetAll() ([]types.User, error) {
	var users []types.User
	err := us.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("error in UserService.GetAll: users bucket does not exist")
		}

		c := bucket.Cursor()
		for key, userRaw := c.First(); key != nil; key, userRaw = c.Next() {
			var user types.User
			err := gob.NewDecoder(bytes.NewBuffer(userRaw)).Decode(&user)
			if err != nil {
				return errors.Wrap(err, "error in UserService.GetAll: problem to decoding the user struct")
			}

			// Strip password for security
			user.Password = []byte{}

			users = append(users, user)
		}

		return nil
	})

	return users, err
}

// getUser is an internal function for retriveing a user from the database
func (us UserService) getUser(email string) (types.User, error) {
	var user types.User
	err := us.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("error in UserService.getUser: users bucket does not exist")
		}

		userRaw := bucket.Get([]byte(email))
		if userRaw == nil {
			fmt.Println("NULL")
			return nil // TODO: Custom Exported Error
		}

		err := gob.NewDecoder(bytes.NewBuffer(userRaw)).Decode(&user)

		return err
	})

	return user, err
}

// Get returns a user by their email
func (us UserService) Get(email string) (types.User, error) {
	user, err := us.getUser(email)

	// Strip password for security
	user.Password = []byte{}

	return user, err
}

// CreateOrEdit creates or edits an existing user if one exists
func (us UserService) CreateOrEdit(email string, user types.User) error {
	// Encode User
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(user); err != nil {
		return errors.Wrap(err, "error in UserService.CreateOrEdit: problem to encoding user struct")
	}
	userRaw := buf.Bytes()

	// Store to DB
	err := us.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("error in UserService.CreateOrEdit: users bucket does not exist")
		}

		err := bucket.Put([]byte(email), userRaw)
		return err
	})

	return err
}

// VerifyLogin takes in a users email & password and checks if they match the users hashed password
func (us UserService) VerifyLogin(email string, password string) (bool, error) {
	// Get User
	user, err := us.getUser(email)
	if err != nil {
		return false, err
	}

	// Compare password against hash
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, err
}

// HasPermission takes a users email & a permission and checks if the user contains that permission.
// It also handles wildcard permissions that can be given to a user.
func (us UserService) HasPermission(email string, permission string) (bool, error) {
	// Get User
	user, err := us.getUser(email)
	if err != nil {
		return false, err
	}

	// Check if user contains permission or equivalent
	for _, permission := range user.Permissions {
		if permission == "*" {
			return true, nil
		} else if permission == permission {
			return true, nil
		}
	}

	return false, nil
}

// HashPassword returns a password in hashed form ready to be stored into the DB
func (us UserService) HashPassword(password []byte) (types.RawPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 11)
	return types.RawPassword(hashedPassword), err
}

// NewUserService creates and initialises a new UserService from a DB connection
func NewUserService(db *bolt.DB) (UserService, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(usersBucket)
		return err
	})

	return UserService{
		db,
	}, err
}
