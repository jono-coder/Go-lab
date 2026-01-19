package etag

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleConditionalGet(w http.ResponseWriter, r *http.Request, version *time.Time) bool {
	etag := MakeWeakETag(version)

	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	w.Header().Set("ETag", etag)
	return false
}

func MakeWeakETag(t *time.Time) string {
	if t == nil {
		return `W/"0"`
	}
	return `W/"` + strconv.FormatInt(t.UnixMicro(), 10) + `"`
}

func ParseETag(h string) (*time.Time, error) {
	if h == "" {
		return nil, fmt.Errorf("missing ETag")
	}

	// Strip weak validator prefix if present
	if strings.HasPrefix(h, "W/") {
		h = strings.TrimPrefix(h, "W/")
	}

	v := strings.Trim(h, `"`)

	if v == "0" {
		return nil, nil // represents NULL updated_at
	}

	micro, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return &time.Time{}, err
	}

	t := time.UnixMicro(micro)
	return &t, nil
}
