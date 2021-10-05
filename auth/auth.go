package auth

import (
	"fmt"

	"github.com/Nerzal/gocloak/v8"
	"github.com/mjehanno/go-ldenerd-api/appconfig/conf"
)

var keyClient gocloak.GoCloak

func GetClient() gocloak.GoCloak {
	if keyClient == nil {
		keyClient = gocloak.NewClient(conf.CurrentConf.KeycloakHost)
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

func (j Jwt) String() string {
	return fmt.Sprintf("{\nAccessToken: %v, \nRefreshToken: %v, \nExpiresIn: %v, \nRefreshExpiresIn: %v, \nTokenType: %v \n}", j.AccessToken, j.RefreshToken, j.ExpiresIn, j.RefreshExpiresIn, j.TokenType)
}
