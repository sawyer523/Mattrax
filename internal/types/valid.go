package types

import "regexp"

// ValidEmail is a regex used to verify an email is valid
var ValidEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// ValidPassword is a regex used to verify a Password is valid
var ValidPassword = regexp.MustCompile(`^[a-zA-Z0-9:\-@ !#$^&*().,?]+$`)
