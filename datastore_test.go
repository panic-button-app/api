package main

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/google/go-cmp/cmp"
)

var fsclient *RemoteFirestoreClient

func TestMain(m *testing.M) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "panic-button-112233")
	if err != nil {
		log.Fatal(err)
	}

	fsclient = NewRemoteFirestoreClient(client).(*RemoteFirestoreClient)
	fsclient.userCollection = client.Collection("Users_Test")

	os.Exit(m.Run())
}

func TestPutGetUser(t *testing.T) {
	user := &User{
		Sub:  "hello@world.com",
		Name: "Michael Scott",
		Contacts: []Contact{
			Contact{
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
