package main

// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
