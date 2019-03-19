package main

import "cloud.google.com/go/datastore"

// User defines what a user looks like in the datastore.
type User struct {
	Key      *datastore.Key `datastore:"__key__"`
	Sub      string         // This is the authentication sub from the oauth provider.
	Name     string
	Contacts []Contact
}
