package models

import "github.com/coreos/go-oidc"

type LoginResponse struct {
	IDToken    oidc.IDToken `json:"idToken"`
	RawIDToken string       `json:"rawIdToken"`
	Email      string       `json:"email"`
}
