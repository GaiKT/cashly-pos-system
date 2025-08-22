package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// OAuthProvider represents different OAuth providers
type OAuthProvider string

const (
	ProviderGoogle   OAuthProvider = "google"
	ProviderFacebook OAuthProvider = "facebook"
)

// OAuthConfig holds OAuth configuration for a provider
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuthUser represents user data from OAuth providers
type OAuthUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Provider      string `json:"provider"`
	ProviderID    string `json:"provider_id"`
	EmailVerified bool   `json:"email_verified"`
}

// OAuthManager handles OAuth operations
type OAuthManager struct {
	configs map[OAuthProvider]*OAuthConfig
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager() *OAuthManager {
	return &OAuthManager{
		configs: make(map[OAuthProvider]*OAuthConfig),
	}
}

// AddProvider adds a new OAuth provider configuration
func (om *OAuthManager) AddProvider(provider OAuthProvider, config *OAuthConfig) {
	om.configs[provider] = config
}

// GetAuthURL generates the OAuth authorization URL
func (om *OAuthManager) GetAuthURL(provider OAuthProvider, state string) (string, error) {
	config, exists := om.configs[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	var baseURL string
	var scopes string

	switch provider {
	case ProviderGoogle:
		baseURL = "https://accounts.google.com/o/oauth2/v2/auth"
		if len(config.Scopes) == 0 {
			scopes = "openid email profile"
		} else {
			scopes = fmt.Sprintf("%v", config.Scopes)
		}
	case ProviderFacebook:
		baseURL = "https://www.facebook.com/v18.0/dialog/oauth"
		if len(config.Scopes) == 0 {
			scopes = "email"
		} else {
			scopes = fmt.Sprintf("%v", config.Scopes)
		}
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	params := url.Values{}
	params.Add("client_id", config.ClientID)
	params.Add("redirect_uri", config.RedirectURL)
	params.Add("scope", scopes)
	params.Add("response_type", "code")
	params.Add("state", state)

	if provider == ProviderGoogle {
		params.Add("access_type", "offline")
		params.Add("prompt", "consent")
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (om *OAuthManager) ExchangeCodeForToken(provider OAuthProvider, code string) (string, error) {
	config, exists := om.configs[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	var tokenURL string
	switch provider {
	case ProviderGoogle:
		tokenURL = "https://oauth2.googleapis.com/token"
	case ProviderFacebook:
		tokenURL = "https://graph.facebook.com/v18.0/oauth/access_token"
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", config.RedirectURL)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", errors.New("access token not found in response")
	}

	return accessToken, nil
}

// GetUserInfo fetches user information using the access token
func (om *OAuthManager) GetUserInfo(provider OAuthProvider, accessToken string) (*OAuthUser, error) {
	var userInfoURL string

	switch provider {
	case ProviderGoogle:
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case ProviderFacebook:
		userInfoURL = "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return om.parseUserInfo(provider, userInfo), nil
}

// parseUserInfo converts provider-specific user info to standard format
func (om *OAuthManager) parseUserInfo(provider OAuthProvider, userInfo map[string]interface{}) *OAuthUser {
	user := &OAuthUser{
		Provider:   string(provider),
		ProviderID: getString(userInfo, "id"),
	}

	switch provider {
	case ProviderGoogle:
		user.Email = getString(userInfo, "email")
		user.Name = getString(userInfo, "name")
		user.FirstName = getString(userInfo, "given_name")
		user.LastName = getString(userInfo, "family_name")
		user.Picture = getString(userInfo, "picture")
		user.EmailVerified = getBool(userInfo, "verified_email")

	case ProviderFacebook:
		user.Email = getString(userInfo, "email")
		user.Name = getString(userInfo, "name")
		user.FirstName = getString(userInfo, "first_name")
		user.LastName = getString(userInfo, "last_name")

		// Facebook picture is nested
		if picture, ok := userInfo["picture"].(map[string]interface{}); ok {
			if data, ok := picture["data"].(map[string]interface{}); ok {
				user.Picture = getString(data, "url")
			}
		}
		user.EmailVerified = true // Facebook email is always verified if provided
	}

	user.ID = user.ProviderID // Set ID to provider ID initially

	return user
}

// Helper functions
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

// GetProviders returns list of configured providers
func (om *OAuthManager) GetProviders() []string {
	providers := make([]string, 0, len(om.configs))
	for provider := range om.configs {
		providers = append(providers, string(provider))
	}
	return providers
}

// IsProviderConfigured checks if a provider is configured
func (om *OAuthManager) IsProviderConfigured(provider OAuthProvider) bool {
	_, exists := om.configs[provider]
	return exists
}
