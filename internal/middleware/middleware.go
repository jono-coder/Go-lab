package middleware

import (
	"net/http"

	"github.com/unrolled/secure"
)

func SecureHandler(next http.Handler) http.Handler {
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:         []string{".*"},
		AllowedHostsAreRegex: true,
		HostsProxyHeaders:    []string{"X-Forwarded-Host"},
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		//ContentSecurityPolicy: "script-src $NONCE",
	})

	return secureMiddleware.Handler(next)
}
