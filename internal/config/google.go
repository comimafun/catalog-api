package internal_config

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config interface {
	googleOauthInit() *oauth2.Config
	exchange(code string) (*oauth2.Token, *domain.Error)
	getUserInfoFromGoogle(token string) (*auth_dto.GoogleUserData, *domain.Error)
	ParseCodeToUserData(code string) (*auth_dto.GoogleUserData, *domain.Error)
	AuthCodeURL() string
}

type config struct{}

// AuthCodeURL implements Config.
func (c *config) AuthCodeURL() string {
	return c.googleOauthInit().AuthCodeURL(os.Getenv("GOOGLE_STATE"), oauth2.AccessTypeOffline)
}

// ParseCodeToUserData implements Config.
func (c *config) ParseCodeToUserData(code string) (*auth_dto.GoogleUserData, *domain.Error) {
	token, err := c.exchange(code)
	if err != nil {
		return nil, err
	}

	userData, err := c.getUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

// exchange implements Config.
func (c *config) exchange(code string) (*oauth2.Token, *domain.Error) {
	if code == "" {
		return nil, domain.NewError(400, errors.New("CODE_IS_EMPTY"), nil)
	}

	token, err := c.googleOauthInit().Exchange(context.Background(), code)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return token, nil
}

// GetUserInfoFromGoogle implements Config.
func (c *config) getUserInfoFromGoogle(token string) (*auth_dto.GoogleUserData, *domain.Error) {
	if token == "" {
		return nil, domain.NewError(400, errors.New("TOKEN_IS_EMPTY"), nil)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	var googleUserData auth_dto.GoogleUserData
	unmarshalErr := json.Unmarshal(userData, &googleUserData)
	if unmarshalErr != nil {
		return nil, domain.NewError(500, unmarshalErr, nil)
	}

	return &googleUserData, nil
}

// GoogleOauthInit implements Config.
func (c *config) googleOauthInit() *oauth2.Config {
	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redicrectURL := os.Getenv("GOOGLE_REDIRECT_URL_DOMAIN")

	return &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redicrectURL,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}
}

func NewConfig() Config {
	return &config{}
}
