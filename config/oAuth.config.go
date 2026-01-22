package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	RedirectURL:  "postmessage", // IMPORTANT for auth-code via JS
	Scopes: []string{
		"openid",
		"email",
		"profile",
	},
	Endpoint: google.Endpoint,
}
