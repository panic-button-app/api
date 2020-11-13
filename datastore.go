package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// UserCollection is what collection is associated with users in firestore.
const UserCollection = "User"

type FirestoreClient interface {
	GetUser(ctx context.Context, user *User) (*User, error)
	PutUser(ctx context.Context, user *User) error
}

// RemoteFirestoreClient
type RemoteFirestoreClient struct {
	fsClient       *firestore.Client
	userCollection *firestore.CollectionRef
}

func NewRemoteFirestoreClient(client *firestore.Client) FirestoreClient {
	c := &RemoteFirestoreClient{fsClient: client}
	c.userCollection = client.Collection(UserCollection)
	return c
}

// GetUser returns the user with the given sub (user.Sub must be populated).
func (rds *RemoteFirestoreClient) GetUser(ctx context.Context, user *User) (*User, error) {
	doc, err := rds.userCollection.Doc(user.GenerateKey()).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by sub: %v", err)
	}
	if err := doc.DataTo(user); err != nil {
		return nil, fmt.Errorf("error while reading user data: %v", err)
	}
	return user, nil
}

// PutUser performs an upsert on the given user.
func (rds *RemoteFirestoreClient) PutUser(ctx context.Context, user *User) error {
	_, err := rds.userCollection.Doc(user.GenerateKey()).Set(ctx, user)
	return err
}
