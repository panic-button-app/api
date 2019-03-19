package main

// Copyright 2019 Google LLC
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

	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/appengine"
)

var (
	firebaseClient *auth.Client
)

func main() {
	// Setup clients.
	app, err := firebase.NewApp(context.Background(), nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	firebaseClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/signIn", signIn)
	http.HandleFunc("/contacts", contacts)
	http.HandleFunc("/press-button", pressButton)
	http.HandleFunc("/ping", ping)
	appengine.Main()
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
}
