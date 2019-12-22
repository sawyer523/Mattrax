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

var usersBucket = []byte("users")

type UserService struct {
	db *bolt.DB
}

func (us UserService) GetUsers() ([]types.User, error) {
	var users []types.User
	err := us.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("error in GetUsers: users bucket does not exist")
		}

		c := bucket.Cursor()
		for key, userRaw := c.First(); key != nil; key, userRaw = c.Next() {
			var user types.User
			err := gob.NewDecoder(bytes.NewBuffer(userRaw)).Decode(&user)
			if err != nil {
				return errors.Wrap(err, "error in GetUsers: problem to decoding the user struct")
			}

			// Strip password for security
			user.Password = []byte{}

			users = append(users, user)
		}

		return nil
	})

	return users, err
}

func (us UserService) getUser(email string) (types.User, error) {
	var user types.User
	err := us.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("error in GetUser: users bucket does not exist")
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

func (us UserService) GetUser(email string) (types.User, error) {
	user, err := us.getUser(email)

	// Strip password for security
	user.Password = []byte{}

	return user, err
}

func (us UserService) CreateOrEditUser(email string, user types.User) error {
	// Encode User
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(user); err != nil {
		return errors.Wrap(err, "error in CreateOrEditUser: problem to encoding user struct")
	}
	userRaw := buf.Bytes()

	// Store to DB
	err := us.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(usersBucket)
		if bucket == nil {
			return errors.New("error in CreateOrEditUser: users bucket does not exist")
		}

		err := bucket.Put([]byte(email), userRaw)
		return err
	})

	return err
}

func (us UserService) VerifyLogin(email string, password string) (bool, error) {
	user, err := us.getUser(email)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}

func (us UserService) HashPassword(password string) (types.RawPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	return types.RawPassword(hashedPassword), err
}

func NewUserService(db *bolt.DB) (UserService, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(usersBucket)
		return err
	})

	return UserService{
		db,
	}, err
}
