package endpoints

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/viddem/huego/internal/utilities"
	"golang.org/x/oauth2"
)

var (
	config       *utilities.HueConfig
	oidcProvider *oidc.Provider
	oauth2Config *oauth2.Config
)

func Init(conf *utilities.HueConfig, provider *oidc.Provider, oauthConfig *oauth2.Config) {
	config = conf
	oidcProvider = provider
	oauth2Config = oauthConfig
}

type ErrorResponse struct {
	Message string `json:"message"`
}
