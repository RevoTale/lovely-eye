package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lovely-eye/server/internal/middleware"
	"github.com/lovely-eye/server/internal/services"
)

type SiteHandler struct {
	siteService *services.SiteService
}

func NewSiteHandler(siteService *services.SiteService) *SiteHandler {
	return &SiteHandler{siteService: siteService}
}

type createSiteRequest struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type updateSiteRequest struct {
	Name string `json:"name"`
}

type siteResponse struct {
	ID        int64  `json:"id"`
	Domain    string `json:"domain"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func (h *SiteHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req createSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Domain == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "Domain and name are required")
		return
	}

	site, err := h.siteService.Create(r.Context(), services.CreateSiteInput{
		Domain: req.Domain,
		Name:   req.Name,
		UserID: claims.UserID,
	})

	if err != nil {
		switch err {
		case services.ErrSiteExists:
			respondError(w, http.StatusConflict, "Site with this domain already exists")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to create site")
		}
		return
	}

	respondJSON(w, http.StatusCreated, siteResponse{
		ID:        site.ID,
		Domain:    site.Domain,
		Name:      site.Name,
		PublicKey: site.PublicKey,
	})
}

func (h *SiteHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sites, err := h.siteService.GetUserSites(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get sites")
		return
	}

	var response []siteResponse
	for _, site := range sites {
		response = append(response, siteResponse{
			ID:        site.ID,
			Domain:    site.Domain,
			Name:      site.Name,
			PublicKey: site.PublicKey,
		})
	}

	if response == nil {
		response = []siteResponse{}
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *SiteHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	siteID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	site, err := h.siteService.GetByID(r.Context(), siteID, claims.UserID)
	if err != nil {
		switch err {
		case services.ErrSiteNotFound:
			respondError(w, http.StatusNotFound, "Site not found")
		case services.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, "Not authorized")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to get site")
		}
		return
	}

	respondJSON(w, http.StatusOK, siteResponse{
		ID:        site.ID,
		Domain:    site.Domain,
		Name:      site.Name,
		PublicKey: site.PublicKey,
	})
}

func (h *SiteHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	siteID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	var req updateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	site, err := h.siteService.Update(r.Context(), siteID, claims.UserID, req.Name)
	if err != nil {
		switch err {
		case services.ErrSiteNotFound:
			respondError(w, http.StatusNotFound, "Site not found")
		case services.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, "Not authorized")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to update site")
		}
		return
	}

	respondJSON(w, http.StatusOK, siteResponse{
		ID:        site.ID,
		Domain:    site.Domain,
		Name:      site.Name,
		PublicKey: site.PublicKey,
	})
}

func (h *SiteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	siteID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	if err := h.siteService.Delete(r.Context(), siteID, claims.UserID); err != nil {
		switch err {
		case services.ErrSiteNotFound:
			respondError(w, http.StatusNotFound, "Site not found")
		case services.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, "Not authorized")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to delete site")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SiteHandler) RegenerateKey(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	siteID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	site, err := h.siteService.RegeneratePublicKey(r.Context(), siteID, claims.UserID)
	if err != nil {
		switch err {
		case services.ErrSiteNotFound:
			respondError(w, http.StatusNotFound, "Site not found")
		case services.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, "Not authorized")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to regenerate key")
		}
		return
	}

	respondJSON(w, http.StatusOK, siteResponse{
		ID:        site.ID,
		Domain:    site.Domain,
		Name:      site.Name,
		PublicKey: site.PublicKey,
	})
}
