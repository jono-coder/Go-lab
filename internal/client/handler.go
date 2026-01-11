package client

import (
	"Go-lab/config"
	"Go-lab/internal/utils/httpconst"
	"Go-lab/internal/utils/session"
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
	return &Handler{
		service: service,
		ctx:     ctx,
		cfg:     cfg,
	}
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	clients, err := h.service.FindAll(h.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(clients) == 0 {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, ToDTOs(clients))
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(h.ctx, h.cfg.TimeoutInSeconds)
	defer cancel()

	// mock //
	ctx = session.ContextWithUserID(ctx, 23174)
	// mock //

	client, err := h.service.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, "client not found")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ToDTO(client))
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Create method not yet implemented")
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Update method not yet implemented")
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Delete method not yet implemented")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set(httpconst.HeaderContentType, httpconst.ContentTypeJson)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println(err)
	}
}
