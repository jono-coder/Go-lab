package security

import (
	"Go-lab/config"
	"Go-lab/internal/utils/httpconst"
	"context"
	"log"
	"net/http"

	"github.com/go-http-utils/headers"
)

type Handler struct {
	config config.AppConfig
	ctx    context.Context
}

func NewHandler(ctx context.Context, config config.AppConfig) *Handler {
	return &Handler{
		config: config,
		ctx:    ctx,
	}
}

// Auth writes a fake access token
func (h *Handler) Auth(w http.ResponseWriter, _ *http.Request) {
	// Fast path for dev
	if h.config.IsDev() {
		log.Println("Auth Test")

		w.Header().Set(headers.ContentType, httpconst.ApplicationJSON)
		w.WriteHeader(http.StatusOK) // explicit status
		w.Write([]byte(`{
		  			"access_token": "fake-test-token-abc123",
		  			"token_type": "Bearer",
		  			"expires_in": 3600
				}`))

		return
	}

	log.Println("Auth for real")
}
