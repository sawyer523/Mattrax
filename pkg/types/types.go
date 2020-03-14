package types

import "regexp"

// IsDNSNameRegex is used to verify if a string is a DNS Name (Domain Name)
var IsDNSNameRegex = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)

// GenericStringRegex is a regex used to verify a simple string
var GenericStringRegex = regexp.MustCompile(`^[a-zA-Z0-9- '"]+$`)

// ValidEmail is a regex used to verify an email is valid
var ValidEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
