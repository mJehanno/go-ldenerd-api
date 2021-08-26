package auth

import "github.com/Nerzal/gocloak/v8"

var keyClient gocloak.GoCloak

func GetClient() gocloak.GoCloak {
	if keyClient == nil {
		keyClient = gocloak.NewClient("http://localhost/")
	}
	return keyClient
}

type User struct {
	Username string
	Password string
}

type Jwt struct {
	AccessToken      string `json:"access_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
}
