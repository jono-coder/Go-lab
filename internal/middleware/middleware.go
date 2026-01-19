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
	})

	return secureMiddleware.Handler(next)
}

func CacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, max-age=30") // Cache per user
		w.Header().Add("Vary", "Authorization")

		next.ServeHTTP(w, r)
	})
}
