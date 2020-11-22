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
	"errors"
	"net/http"
	"os"
	"strings"

	pberrors "github.com/panic-button-app/api/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var auth0Domain = os.Getenv("AUTH0_DOMAIN")

const sessionName = "user-session"

// HandlerWithUser defines what a http handler that requires a user looks like.
type HandlerWithUser = func(w http.ResponseWriter, r *http.Request, user *User)

// SessionMiddleware is a http middleware that ensures a user is available and verified
// before calling the next handler.
func SessionMiddleware(next HandlerWithUser, fs FirestoreClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Ensure that a session already exists.
		if session.IsNew {
			http.Error(w, "session does not exist", http.StatusUnauthorized)
			return
		}

		val := session.Values["UserKey"]
		if _, ok := val.(string); !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userKey := val.(string)

		user, err := fs.GetUser(r.Context(), &User{Sub: userKey})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next(w, r, user)
	}
}

// SaveUserSession saves the given user to the session.
func SaveUserSession(w http.ResponseWriter, r *http.Request, user *User) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values["UserKey"] = user.Sub
	return session.Save(r, w)
}

func ExtractAuth0JWT(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	tokenString := authHeaderParts[1]

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getPemCert(r.Context(), token)
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid certificate")
	}

	return claims.Subject, nil
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func getPemCert(ctx context.Context, token *jwt.Token) (string, error) {
	cert := ""
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+auth0Domain+"/.well-known/jwks.json", nil)
	if err != nil {
		return cert, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return cert, pberrors.Annotate(err, pberrors.CodeDependentServiceFailure)
	}
	defer resp.Body.Close()

	var jwks = &Jwks{}
	err = json.NewDecoder(resp.Body).Decode(jwks)
	if err != nil {
		return cert, err
	}

	for _, k := range jwks.Keys {
		if token.Header["kid"] == k.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + k.X5c[0] + "\n-----END CERTIFICATE-----"
			break
		}
	}

	if cert == "" {
		return cert, pberrors.Annotate(
			errors.New("unable to find appropriate key"), pberrors.CodeUnauthorized)
	}

	return cert, nil
}
