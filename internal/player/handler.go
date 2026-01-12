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
	ctx     context.Context
	cfg     config.AppConfig
}

func NewHandler(ctx context.Context, service *Service, cfg config.AppConfig) *Handler {
	if err := validate.Required("ctx", ctx); err != nil {
		panic(err)
	}
	if err := validate.Required("service", service); err != nil {
		panic(err)
	}
	return &Handler{
		service: service,
		ctx:     ctx,
		cfg:     cfg,
	}
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	players, err := h.service.FindAll(h.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(players) == 0 {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, ToDTOs(players))
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(h.ctx, h.cfg.TimeoutInSeconds)
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
	resourceId := chi.URLParam(r, "resource_id")
	if resourceId == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(h.ctx, h.cfg.TimeoutInSeconds)
	defer cancel()

	player, err := h.service.FindByResourceId(ctx, resourceId)
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
		w.WriteHeader(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(h.ctx, h.cfg.TimeoutInSeconds)
	defer cancel()

	player, err := h.service.Checkin(ctx, uint(id))
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set(httpconst.HeaderContentType, httpconst.ContentTypeJson)
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Println(err)
	}
}
