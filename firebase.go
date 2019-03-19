package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

// ErrMissingToken means that an ID token could not be found in
// the request. Usually set in the "Authorization" header.
var ErrMissingToken = errors.New("Missing ID Token in request")

const validAuthType = "Bearer"

// authenticate takes a http request, and either returns the string
// of the sub for the user, or an error.
func authenticate(ctx context.Context, r *http.Request, client *auth.Client) (string, error) {
	// Extract the id token from the request.
	authHeader := r.Header.Get("Authorization")
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) < 2 || headerParts[0] != validAuthType {
		return "", ErrMissingToken
	}

	token, err := client.VerifyIDToken(ctx, headerParts[1])
	if err != nil {
		return "", err
	}

	return token.Subject, nil
}
