package api

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/viddem/huego/internal/api/endpoints"
	"github.com/viddem/huego/internal/utilities"
	"golang.org/x/oauth2"
)

var (
	config       *utilities.HueConfig
	oidcProvider *oidc.Provider
	oauth2Config *oauth2.Config
)

func Init(conf *utilities.HueConfig) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, conf.OIDCIssuer)
	if err != nil {
		log.Fatalf("Failed to initialize OIDC provider: %v", err)
	}
	oidcProvider = provider

	oauth2Config = &oauth2.Config{
		ClientID:     conf.OIDCClientID,
		ClientSecret: conf.OIDCClientSecret,
		RedirectURL:  conf.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	router := gin.Default()
	store := cookie.NewStore([]byte(conf.Secret))
	router.Use(sessions.Sessions("auth", store))
	endpoints.Init(conf, oidcProvider, oauth2Config)
	config = conf

	v1 := router.Group("/api/")
	{
		auth := v1.Group("")
		auth.Use(CheckAuth())
		{
			auth.GET("/lamps", endpoints.GetLamps)
			auth.POST("/lamps", endpoints.SetLamps)
			auth.POST("/lamps/:id", endpoints.SetLamp)
		}
		v1.POST("/auth", endpoints.Auth)
		v1.POST("/logout", endpoints.Logout)
	}

	err = router.Run()
	if err != nil {
		log.Fatalf("Failed to start webserver due to err: %s\n", err)
	}
}

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		idToken := session.Get("id_token")

		if idToken == nil {
			InitializeAuth(c)
			return
		}

		verifier := oidcProvider.Verifier(&oidc.Config{ClientID: config.OIDCClientID})
		_, err := verifier.Verify(context.Background(), idToken.(string))
		if err != nil {
			log.Printf("Invalid token: %v", err)
			InitializeAuth(c)
			return
		}
	}
}

func InitializeAuth(c *gin.Context) {
	state := generateRandomState()
	session := sessions.Default(c)
	session.Set("oauth_state", state)
	session.Save()

	authURL := oauth2Config.AuthCodeURL(state)

	c.Header("location", authURL)
	c.String(http.StatusUnauthorized, authURL)
	c.Abort()
}

func generateRandomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(65 + (i % 26)) // Simple random string
	}
	return string(b)
}
