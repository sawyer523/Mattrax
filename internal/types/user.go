package types

import (
	"errors"
	"regexp"
	"time"
)

// RawPassword is a string that contains the hashed + salted password.
// This type should NEVER EVER EVER contain a plain text password!!!!!!
// It is not exposed outside this package to force the use of helpers to keep
// the password secure!
type RawPassword string

// Permissions contains a list of allowed permissions for a user.
type Permissions []string

// An Action is an event taken by a user.
type Action struct {
	EventName        string
	EventDescription string
	Time             time.Time
}

// A User is either an administrator or end user who owns a managed device.
type User struct {
	DisplayName string
	Email       string
	Password    RawPassword
	Activity    []Action    `graphql:",optional"` // TODO: API Read only
	Permissions Permissions `graphql:",optional"` // TODO: API Require PERM to change
}

// Regex's are used to verify the users input
var SafeString = regexp.MustCompile(`^[a-zA-Z0-9:\-@ !#$^&*().,?]+$`)

// ValidEmail is a regex used to verify an email is valid
var ValidEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// Verify checks that the structs fields are valid
func (user User) Verify() error {
	if user.DisplayName != "" && !SafeString.MatchString(user.DisplayName) {
		return errors.New("invalid user: invalid DisplayName '" + user.DisplayName + "'")
	}

	if user.Email != "" && !ValidEmail.MatchString(user.Email) {
		return errors.New("invalid user: invalid Email '" + user.Email + "'")
	}

	if user.Password == "" {
		return errors.New("invalid user: missing Password")
	}

	return nil
}

// ErrUserNotFound is the error returned if a user can't be found
var ErrUserNotFound = errors.New("Error: User not found")

// UserService contains the implemented functionality for users
type UserService interface {
	GetAll() ([]User, error)
	Get(email string) (User, error)
	CreateOrEdit(email string, user User) error
	VerifyLogin(email string, password string) (bool, error)
	HasPermission(email string, permission string) (bool, error)
	HashPassword(password []byte) (RawPassword, error)
}
