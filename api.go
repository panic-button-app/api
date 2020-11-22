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
	"encoding/json"
	"net/http"
	"os"

	"github.com/panic-button-app/api/errors"

	"log"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/handlers"
)

var (
	firestoreClient FirestoreClient
)

func main() {
	// Setup clients.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		log.Fatal(err)
	}
	firestoreClient = NewRemoteFirestoreClient(client)

	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r := http.NewServeMux()

	r.HandleFunc("/signIn", signIn)
	r.HandleFunc("/getUser", SessionMiddleware(getUser, firestoreClient))
	r.HandleFunc("/pressButton", pressButton)
	r.HandleFunc("/ping", ping)

	// Wrap our server with our gzip handler to gzip compress all responses.
	http.ListenAndServe(":"+port,
		handlers.RecoveryHandler()(handlers.CompressHandler(r)))
}

// ErrorMessage defines how the API returns error messages back
// to the client.
type ErrorMessage struct {
	Code    int
	Message string
}

func handleError(err error, w http.ResponseWriter) {
	// Check if the error has been annotated with the errors package.
	annotated, ok := err.(errors.Error)
	if !ok {
		// Default to internal.
		annotated = errors.Annotate(err, errors.CodeInternal).(errors.Error)
	}

	log.Println(annotated.Error())

	status := errors.HTTPMapping[annotated.Code]
	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorMessage{
		Code: status,
	}); err != nil {
		log.Println(err)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	// Get the auth0 token from the request.
	auth0Token := ""
	_ = auth0Token

	userSub, err := ExtractAuth0JWT(r)
	if err != nil {
		handleError(err, w)
		return
	}
	log.Printf("userSub: %v\n", userSub)

}

// Retrieves the user.
func getUser(w http.ResponseWriter, r *http.Request, user *User) {
	if err := json.NewEncoder(w).Encode(user); err != nil {
		handleError(err, w)
		return
	}
}

func pressButton(w http.ResponseWriter, r *http.Request) {

}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
