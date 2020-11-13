package main

import (
	"crypto/sha256"
	"fmt"
)

// User defines what a user looks like in the datastore.
type User struct {
	Sub string // This is the authentication sub from the oauth
	// provider. Typically an email address.
	Name     string
	Contacts []Contact
}

// GenerateKey returns a unique deterministic hash for this user.
func (u *User) GenerateKey() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(u.Sub)))
}
