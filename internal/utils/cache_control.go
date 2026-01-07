package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	equals                  = "="
	headerCacheControl      = "Cache-Control"
	headerExpires           = "Expires"
	dirMaxAge               = "max-age"
	dirStaleWhileRevalidate = "stale-while-revalidate"
	dirMustRevalidate       = "must-revalidate"
	dirPublic               = "public"
	dirPrivate              = "private"
	separator               = ", "
)

var (
	Low     = CacheControl(10*time.Second, 30*time.Second, false)
	Medium  = CacheControl(5*time.Minute, 10*time.Minute, false)
	High    = CacheControl(1*time.Hour, 2*time.Hour, false)
	Highest = CacheControl(12*time.Hour, 24*time.Hour, false)
)

func CacheControl(maxAge, staleWhileRevalidate time.Duration, public bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			directives := make([]string, 0, 4)

			if maxAge > 0 {
				directives = append(directives, dirMaxAge+equals+fmt.Sprint(int(maxAge.Seconds())))
			}
			if staleWhileRevalidate > 0 {
				directives = append(directives, dirStaleWhileRevalidate+equals+fmt.Sprint(int(staleWhileRevalidate.Seconds())))
			}

			directives = append(directives, dirMustRevalidate)
			if public {
				directives = append(directives, dirPublic)
			} else {
				directives = append(directives, dirPrivate)
			}

			w.Header().Set(headerCacheControl, strings.Join(directives, separator))

			if maxAge > 0 {
				w.Header().Set(headerExpires, time.Now().Add(maxAge).UTC().Format(http.TimeFormat))
			}

			next.ServeHTTP(w, r)
		})
	}
}
