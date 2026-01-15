package player

import (
	"Go-lab/config"
	"Go-lab/internal/utils/httpconst"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	cfg     config.AppConfig
}

func NewHandler(ctx context.Context, service *Service, cfg config.AppConfig) *Handler {
	if err := validate.Required("service", service); err != nil {
		panic(err)
	}
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	players, err := h.service.FindAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, ToDTOs(players))
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	player, err := h.service.FindById(ctx, uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, "player not found")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ToDTO(player))
}

func (h Handler) GetResource(w http.ResponseWriter, r *http.Request) {
	resourceID := chi.URLParam(r, "resource_id")
	if resourceID == "" {
		http.Error(w, "invalid resource id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	player, err := h.service.FindByResourceId(ctx, resourceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, "player not found")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ToDTO(player))
}

func (h Handler) Checkin(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	player, err := h.service.Checkin(ctx, uint(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Not a missing route — a state or timing conflict
			writeJSON(w, http.StatusConflict, "player not ready or not found")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ToDTO(player))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set(httpconst.HeaderContentType, httpconst.ContentTypeJson)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println(err)
	}
}
