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
	r.HandleFunc("/contacts", contacts)
	r.HandleFunc("/press-button", pressButton)
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

func sendError(status int, msg string, w http.ResponseWriter) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(ErrorMessage{
		Code:    status,
		Message: msg,
	})
	if err != nil {
		log.Println(err)
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	// Verify the user's identity.
}

// Contact management.
func contacts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:

	}
}

func pressButton(w http.ResponseWriter, r *http.Request) {

}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
