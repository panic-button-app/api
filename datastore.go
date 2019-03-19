package main

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// UserKind is what kind is associated with users in the datastore.
const UserKind = "User"

type DatastoreClient interface {
	GetUserBySub(ctx context.Context, sub string) (*User, error)
	PutUser(ctx context.Context, user *User) error
}

type RemoteDatastoreClient struct {
	dsClient *datastore.Client
}

func NewRemoteDatastoreClient(client *datastore.Client) DatastoreClient {
	return &RemoteDatastoreClient{client}
}

func (rds *RemoteDatastoreClient) GetUserBySub(ctx context.Context, sub string) (*User, error) {
	query := datastore.NewQuery(UserKind).Filter("Sub =", sub).Limit(1)
	t := rds.dsClient.Run(ctx, query)
	var user *User
	for {
		u := new(User)
		_, err := t.Next(u)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		user = u
	}

	return user, nil
}

func (rds *RemoteDatastoreClient) PutUser(ctx context.Context, user *User) error {
	// If the user does not have a key, get one from the server.
	if user.Key == nil {
		user.Key = datastore.IncompleteKey(UserKind, nil)
	}

	_, err := rds.dsClient.Put(ctx, user.Key, user)
	return err
}
