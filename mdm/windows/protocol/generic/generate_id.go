package generic

import uuid "github.com/satori/go.uuid"

// GenerateID returns a new random UUID string to be used as an ID
func GenerateID() string {
	return uuid.NewV4().String()
}
