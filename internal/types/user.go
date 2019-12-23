package types

import "time"

// RawPassword is a string that contains the hashed + salted password.
// This type should NEVER EVER EVER contain a plain text password!!!!!!
// It is not exposed outside this package to force the use of helpers to keep
// the password secure!
type RawPassword []byte

// Permissions contains a list of allowed permissions for a user.
type Permissions []string

// A User is either an administrator or end user who owns a managed device.
type User struct {
	DisplayName string
	Email       string      `sqlgen:",primary"`
	Password    RawPassword `graphql:"-"`
	Activity    []Action    `graphql:",optional"` // TODO: API Read only
	Permissions Permissions `graphql:",optional"` // TODO: API Require PERM to change
}

// An Action is an event taken by a user.
type Action struct {
	EventName        string
	EventDescription string
	Time             time.Time
}

// UserService contains the implemented functionality for users
type UserService interface {
	GetAll() ([]User, error)
	Get(email string) (User, error)
	CreateOrEdit(email string, user User) error
	VerifyLogin(email string, password string) (bool, error)
	HasPermission(email string, permission string) (bool, error)
	HashPassword(password []byte) (RawPassword, error)
}
