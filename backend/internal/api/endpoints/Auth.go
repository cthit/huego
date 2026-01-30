package endpoints

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Code  *string `json:"code"`
	State *string `json:"state"`
}

func Auth(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error: failed to read json data, err: %s\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to read request data",
		})
		return
	}

	var receivedAuth authRequest
	err = json.Unmarshal(jsonData, &receivedAuth)
	if receivedAuth.Code == nil {
		log.Printf("No code in request")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid or missing code",
		})
		return
	}

	session := sessions.Default(c)
	stored_state := session.Get("oauth_state")
	if stored_state == nil || (receivedAuth.State != nil && *receivedAuth.State != stored_state.(string)) {
		log.Printf("Invalid state parameter")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid state parameter",
		})
		return
	}

	ctx := context.Background()

	oauth2Token, err := oauth2Config.Exchange(ctx, *receivedAuth.Code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Failed to exchange authorization code",
		})
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Printf("No id_token field in oauth2 token")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "No ID token received",
		})
		return
	}

	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: config.OIDCClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("Failed to verify ID Token: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid ID token",
		})
		return
	}

	session.Set("id_token", rawIDToken)
	session.Set("access_token", oauth2Token.AccessToken)
	session.Options(sessions.Options{
		MaxAge: int(oauth2Token.Expiry.Unix()),
	})
	err = session.Save()
	if err != nil {
		log.Printf("Failed to create session: %v\n", err)
		c.JSON(500, ErrorResponse{
			Message: "Failed to create session",
		})
		return
	}

	var claims struct {
		Name    string `json:"name"`
		Subject string `json:"sub"`
	}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Failed to parse claims: %v", err)
	}

	c.String(http.StatusOK, "Session created")
}
