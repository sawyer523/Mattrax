package generic

import "regexp"

// validMessageID is a regex used to verify a MessageID is valid
var validMessageID = regexp.MustCompile(`^[a-zA-Z0-9:\-]+$`)

// validBinarySecurityToken is a regex used to verify a base64 encoding binary securiry token is valid
var validBinarySecurityToken = regexp.MustCompile(`^[a-zA-Z0-9+/]+$`)
