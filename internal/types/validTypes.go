package types

import "regexp"

// IsDNSNameRegex is used to verify if a string is a DNS Name (Domain Name)
var IsDNSNameRegex = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)
