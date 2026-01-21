package player

import (
	"Go-lab/config"
	"Go-lab/internal/middleware/etag"
	"Go-lab/internal/utils/httpconst"
	"Go-lab/internal/utils/paging"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	cfg     config.AppConfig
}

func NewHandler(service *Service, cfg config.AppConfig) *Handler {
	if err := validate.Get().Var(service, "required"); err != nil {
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

	page, err := paging.ParsePage(chi.URLParam(r, "page"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := paging.ParseLimit(chi.URLParam(r, "limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	players, err := h.service.FindAll(ctx, paging.NewPaging(page, limit))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dtos, err := ToDTOs(players)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, dtos)
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

	if etag.HandleConditionalGet(w, r, player.UpdatedAt) {
		return
	}

	dto, err := ToDTO(player)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto)
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

	if etag.HandleConditionalGet(w, r, player.UpdatedAt) {
		return
	}

	dto, err := ToDTO(player)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dto)
}

func (h Handler) Checkin(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	header := r.Header.Get("If-Match")
	version, err := etag.ParseETag(header)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	player, err := h.service.Checkin(ctx, uint(id), version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusConflict, "player already modified by another request, please refresh and retry.")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	dto, err := ToDTO(player)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dto)
}

type UpdateDto struct {
	Id          *uint		`db:"id"`
	Name        string		`db:"name"`
	Description *string		`db:"description"`
	UpdatedAt   *time.Time	`db:"updated_at"`
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid player id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	header := r.Header.Get("If-Match")
	version, err := etag.ParseETag(header)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if version == nil {
		http.Error(w, "invalid eTag", http.StatusBadRequest)
		return
	}
	
	var _dto *UpdateDto
	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&_dto); err != nil {
		http.Error(w, "Bad JSON format", http.StatusBadRequest)
		return
	}

	_id := uint(id)
	_dto.Id = &_id

	_, err = h.service.Update(ctx, _dto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusConflict, "player already modified by another request, please refresh and retry.")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.TimeoutInSeconds)
	defer cancel()

	var dto *DTO

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Bad JSON format", http.StatusBadRequest)
		return
	}

	player, err := ToEntity(*dto)

	id, err := h.service.Create(ctx, player)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, "player not found")
			return
		}
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	etag.HandleConditionalGet(w, r, nil)

	writeJSON(w, http.StatusCreated, id)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set(httpconst.HeaderContentType, httpconst.ContentTypeJson)
	w.WriteHeader(status)
	if v == nil{
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println(err)
	}
}
