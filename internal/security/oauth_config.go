package security

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OAuthConfig struct {
	Client *resty.Client
}

type loggingTransport struct {
	base http.RoundTripper
}

func (t loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(">>>>> Sending:", req.Method, req.URL.String())
	fmt.Println("Headers:", req.Header)
	auth := req.Header.Get("Authorization")
	if auth == "" {
		fmt.Println("No Authorization header")
	} else {
		fmt.Println("Authorization header:", auth)
	}
	return t.base.RoundTrip(req)
}

func NewOAuthConfig(ctx context.Context, baseUrl string) *OAuthConfig {
	cfg := &clientcredentials.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		TokenURL:     baseUrl + "/security/oauth/token",
		Scopes:       []string{"api:read", "api:write"},
	}

	c := cfg.Client(ctx)

	logTransport(c)

	Client := resty.NewWithClient(c).
		SetBaseURL(baseUrl).
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	res := &OAuthConfig{
		Client: Client,
	}

	return res
}

func logTransport(c *http.Client) {
	oauthTr, ok := c.Transport.(*oauth2.Transport)
	if !ok {
		log.Fatal("transport is not of type oauth2.Transport")
	}
	base := oauthTr.Base
	if base == nil {
		base = http.DefaultTransport
	}

	// Wrap the base in loggingTransport
	oauthTr.Base = loggingTransport{base: base}
}
