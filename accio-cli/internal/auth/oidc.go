package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IDToken      string    `json:"id_token"`
	Expiry       time.Time `json:"expiry"`
}

const (
	clientID     = "accio-cli"
	clientSecret = "" // Public client in Keycloak
	issuerURL    = "http://localhost:8888/realms/accio"
	callbackURL  = "http://localhost:18080/callback"
	listenAddr   = ":18080"
)

func Login() (*Token, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	state, err := generateRandomState()
	if err != nil {
		return nil, err
	}

	codeCh := make(chan string)
	errCh := make(chan error)

	server := &http.Server{Addr: listenAddr}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "State mismatch", http.StatusBadRequest)
			errCh <- fmt.Errorf("state mismatch")
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			errCh <- fmt.Errorf("code not found")
			return
		}

		fmt.Fprintf(w, "Login successful! You can close this window.")
		codeCh <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %v", err)
		}
	}()

	authURL := oauth2Config.AuthCodeURL(state)
	fmt.Printf("Opening browser for authentication: %s\n", authURL)
	// In a real app we'd use 'pkg/browser' to open this URL
	// exec.Command("xdg-open", authURL).Start()

	select {
	case code := <-codeCh:
		server.Shutdown(ctx)
		oauth2Token, err := oauth2Config.Exchange(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange token: %v", err)
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			return nil, fmt.Errorf("no id_token field in oauth2 token")
		}

		return &Token{
			AccessToken:  oauth2Token.AccessToken,
			RefreshToken: oauth2Token.RefreshToken,
			IDToken:      rawIDToken,
			Expiry:       oauth2Token.Expiry,
		}, nil
	case err := <-errCh:
		server.Shutdown(ctx)
		return nil, err
	case <-time.After(5 * time.Minute):
		server.Shutdown(ctx)
		return nil, fmt.Errorf("login timed out")
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SaveToken(t *Token) error {
	// Simple file storage for dev. In prod use keyring.
	home, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.accio", home)
	os.MkdirAll(configDir, 0700)

	data, _ := json.Marshal(t)
	return os.WriteFile(fmt.Sprintf("%s/token.json", configDir), data, 0600)
}

func LoadToken() (*Token, error) {
	home, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.accio", home)
	data, err := os.ReadFile(fmt.Sprintf("%s/token.json", configDir))
	if err != nil {
		return nil, err
	}
	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
