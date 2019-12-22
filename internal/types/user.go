package types

import "time"

// rawPassword is a string that contains the hashed + salted password.
// This type should NEVER EVER EVER contain a plain text password!!!!!!
// It is not exposed outside this package to force the use of helpers to keep
// the password secure!
type RawPassword []byte

// Permissions contains a list of allowed permissions for a user.
type Permissions []string

// A User is either an administrator or end user who owns a managed device.
type User struct {
	DisplayName string
	Email       string
	Password    RawPassword `graphql:"-"`
	Activity    []Action    // TODO: API Read only
	Permissions Permissions // TODO: API Require PERM to change
}

// An Action is an event taken by a user.
type Action struct {
	EventName        string
	EventDescription string
	Time             time.Time
}

// UserService contains the implemented functionality for users
type UserService interface {
	GetUsers() ([]User, error)
	GetUser(email string) (User, error)
	CreateOrEditUser(email string, user User) error
	VerifyLogin(email string, password string) (bool, error)
	HashPassword(password string) (RawPassword, error)
}
