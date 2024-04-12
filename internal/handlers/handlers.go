package handlers

import (
	"banners/internal/models"
	"banners/internal/usecase/BMS"
	"encoding/json"
	"errors"
	"net/http"
)

// UserBannerHandler handles GET request with
// Query params: tag_id, feature_id, use_last_version;
// Header params: token
func (c *Controller) UserBannerHandler(w http.ResponseWriter, r *http.Request) {}

// BannersSortedHandler handles GET request with
// Query params: tag_id, feature_id, limit, offset;
// Header params: token
func (c *Controller) BannersSortedHandler(w http.ResponseWriter, r *http.Request) {

}

// CreateBannerHandler handles POST request with
// Header params: token
func (c *Controller) CreateBannerHandler(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = c.Usecases.CreateBanner(r.Context(), banner)
	if errors.Is(err, BMS.ErrFeatureTagPairAlreadyExists) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PatchBannerHandler handles PATCH request with
// Query params: id;
// Header params: token
func (c *Controller) PatchBannerHandler(w http.ResponseWriter, r *http.Request) {}

// DeleteBannerHandler handles DELETE request with
// Query params: id;
// Header params: token
func (c *Controller) DeleteBannerHandler(w http.ResponseWriter, r *http.Request) {}
