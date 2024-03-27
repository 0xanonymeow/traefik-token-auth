// Package traefik_token_auth .
package traefik_token_auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Config the plugin configuration.
type Config struct {
	HeaderField  string
	HashedToken  string
	RemoveHeader bool
	Algorithm    string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderField:  "X-Api-Token",
		HashedToken:  "",
		RemoveHeader: true,
		Algorithm:    "sha256",
	}
}

// TokenAuth a traefik_token_auth plugin.
type TokenAuth struct {
	next         http.Handler
	name         string
	headerField  string
	hashedToken  string
	removeHeader bool
	algorithm    string
}

// New created a new token auth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Printf("Creating plugin: %s instance: %+v, ctx: %+v\n", name, *config, ctx)

	return &TokenAuth{
		next:         next,
		name:         name,
		headerField:  config.HeaderField,
		hashedToken:  config.HashedToken,
		removeHeader: config.RemoveHeader,
		algorithm:    config.Algorithm,
	}, nil
}

func (ta *TokenAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if ta.headerField == "" || ta.hashedToken == "" {
		fmt.Printf("header field or hashed token is empty")
		rw.WriteHeader(http.StatusUnauthorized)
	}

	token := req.Header.Get(ta.headerField)

	switch ta.algorithm {
	case "sha256":
		sum := sha256.Sum256([]byte(token))

		if hex.EncodeToString(sum[:]) != ta.hashedToken {
			fmt.Printf("token is invalid")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
	case "bcrypt":
		err := bcrypt.CompareHashAndPassword([]byte(ta.hashedToken), []byte(token))
		if err != nil {
			fmt.Printf("token is invalid")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	if ta.removeHeader {
		req.Header.Del(ta.headerField)
	}

	ta.next.ServeHTTP(rw, req)
}
