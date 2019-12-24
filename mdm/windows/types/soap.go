package wtypes

import uuid "github.com/satori/go.uuid"

// MustUnderstand is a easily way to create SOAP tags with s:mustUnderstand
type MustUnderstand struct {
	MustUnderstand string `xml:"s:mustUnderstand,attr"`
	Value          string `xml:",innerxml"`
}

// GenerateActivityID returns a new random UUID string to be used as an ActivityID
func GenerateActivityID() string {
	return uuid.NewV4().String()
}
