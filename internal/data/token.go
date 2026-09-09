package data

import "time"

type (
	TokenBody struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	RefreshTokenBody struct {
		Refresh string `json:"refresh_token"`
	}

	TokenPayload struct {
		AuthenticationToken struct {
			Token   string `json:"token"`
			Refresh string `json:"refresh_token"`
			Expiry  string `json:"expiry"`
			Scope   string `json:"scope"`
		} `json:"authentication_token"`
	}

	// Cacher interface
	Cacher interface {
		LoadToken(string) (*TokenPayload, error)
		SaveToken(*TokenPayload) (string, error)
		ReplaceToken(string, *TokenPayload) error
		RemoveToken(string) error
	}
)

func (t *TokenPayload) Save(c Cacher) (string, error) {
	return c.SaveToken(t)
}

func (t *TokenPayload) Replace(c Cacher, ID string) error {
	return c.ReplaceToken(ID, t)
}

func (t *TokenPayload) IsExpired() bool {
	exp, err := time.Parse("2006-01-02 15:04:05", t.AuthenticationToken.Expiry)
	if err != nil {
		return false
	}
	return time.Now().Before(exp)
}

func LoadToken(c Cacher, ID string) (*TokenPayload, error) {
	return c.LoadToken(ID)
}

func RemoveToken(c Cacher, ID string) error {
	return c.RemoveToken(ID)
}
