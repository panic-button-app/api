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
	"os"
	"sync"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/google/go-cmp/cmp"
)

var fsclient *RemoteFirestoreClient
var setupOnce = &sync.Once{}

func ensureSetup(t *testing.T) {
	setupOnce.Do(func() {
		ctx := context.Background()
		client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT_ID"))
		if err != nil {
			t.Fatal(err)
		}

		fsclient = NewRemoteFirestoreClient(client).(*RemoteFirestoreClient)
		fsclient.userCollection = client.Collection("Users_Test")
	})
}

func TestPutGetUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	ensureSetup(t)

	user := &User{
		Sub:  "hello@world.com",
		Name: "Michael Scott",
		Contacts: []Contact{
			{
				Name:        "Pam",
				PhoneNumber: "123456",
			},
		},
	}

	ctx := context.Background()
	err := fsclient.PutUser(ctx, user)
	if err != nil {
		t.Errorf("PutUser(ctx, user): Got %q, want nil", err)
	}

	// Test retrieval
	got, err := fsclient.GetUser(ctx, &User{Sub: "hello@world.com"})
	if err != nil {
		t.Errorf("GetUser(ctx, user): Got %q, want nil", err)
	}
	if diff := cmp.Diff(user, got); diff != "" {
		t.Errorf("GetUser(ctx, user) mismatch (-want +got):\n%s", diff)
	}
}

func TestSubKeyGeneration(t *testing.T) {
	tests := []struct {
		name string
		Sub  string
		want string
	}{
		{
			name: "empty string",
			Sub:  "",
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "sub is bob",
			Sub:  "bob",
			want: "81b637d8fcd2c6da6359e6963113a1170de795e4b725b84d1e0b4cfd9ec58ce9",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := &User{Sub: test.Sub}

			// SHA256 of an empty string
			got := user.GenerateKey()
			if got != test.want {
				t.Errorf("user.GenerateKey(): got %v, want %v", got, test.want)
			}
		})
	}
}
