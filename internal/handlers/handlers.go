package handlers

import (
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase"
	"banners/internal/usecase/BMS"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"strconv"
)

// UserBannerHandler handles GET request with
// Query params: tag_id, feature_id, use_last_version;
// Header params: token
func (c *Controller) UserBannerHandler(w http.ResponseWriter, r *http.Request) {
	tagID := r.URL.Query().Get("tag_id")
	tag, err := strconv.Atoi(tagID)
	if tagID == "" {
		ProcessError(w, ErrNoTag, http.StatusBadRequest)
		return
	}
	if err != nil {
		ProcessError(w, ErrIncorrectTag, http.StatusBadRequest)
		return
	}

	featureID := r.URL.Query().Get("feature_id")
	feature, err := strconv.Atoi(featureID)
	if featureID == "" {
		ProcessError(w, ErrIncorrectFeature, http.StatusBadRequest)
		return
	}
	if err != nil {
		ProcessError(w, ErrIncorrectFeature, http.StatusBadRequest)
		return
	}

	banner, err := c.Usecases.UserBanner(r.Context(), int32(tag), int32(feature))
	if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(banner)
	if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
	}
}

// BannersSortedHandler handles GET request with
// Query params: tag_id, feature_id, limit, offset;
// Header params: token
func (c *Controller) BannersSortedHandler(w http.ResponseWriter, r *http.Request) {
	var tagSQLC pgtype.Int4
	tagID := r.URL.Query().Get("tag_id")
	if tagID != "" {
		tag, err := strconv.Atoi(tagID)
		if err != nil {
			ProcessError(w, ErrNoTag, http.StatusBadRequest)
			return
		}
		tagSQLC = pgtype.Int4{
			Int32: int32(tag),
			Valid: true,
		}
	}

	var featureSQLC pgtype.Int4
	featureID := r.URL.Query().Get("feature_id")
	if featureID != "" {
		feature, err := strconv.Atoi(featureID)
		if err != nil {
			ProcessError(w, ErrNoFeature, http.StatusBadRequest)
			return
		}
		featureSQLC = pgtype.Int4{
			Int32: int32(feature),
			Valid: true,
		}
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = strconv.Itoa(DefaultLimit)
	}
	limitVal, err := strconv.Atoi(limit)
	if err != nil {
		ProcessError(w, ErrIncorrectLimit, http.StatusBadRequest)
		return
	}

	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = strconv.Itoa(DefaultOffset)
	}
	offsetVal, err := strconv.Atoi(offset)
	if err != nil {
		ProcessError(w, ErrIncorrectOffset, http.StatusBadRequest)
		return
	}

	bannersList, err := c.Usecases.ListBanners(r.Context(), banners.ListBannersParams{
		TagID:     tagSQLC,
		FeatureID: featureSQLC,
		OffsetVal: int32(offsetVal),
		LimitVal:  int32(limitVal),
	})
	if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
		return
	}
	bannersListDTO := usecase.NewListBannersDTO(bannersList)
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(bannersListDTO)
	if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
		return
	}
}

// CreateBannerHandler handles POST request with
// Header params: token
func (c *Controller) CreateBannerHandler(w http.ResponseWriter, r *http.Request) {
	var banner usecase.BannerJsonDTO
	err := json.NewDecoder(r.Body).Decode(&banner)
	if err != nil {
		ProcessError(w, ErrIncorrectBannerContent, http.StatusBadRequest)
		return
	}

	err = c.Usecases.CreateBanner(r.Context(), banner)
	if errors.Is(err, BMS.ErrFeatureTagPairAlreadyExists) {
		ProcessError(w, err, http.StatusBadRequest)
	} else if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
		return
	}
}

// PatchBannerHandler handles PATCH request with
// Query params: id;
// Header params: token
func (c *Controller) PatchBannerHandler(w http.ResponseWriter, r *http.Request) {
	var update usecase.UpdateBannerDTO
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		ProcessError(w, ErrIncorrectBannerContent, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		ProcessError(w, ErrIncorrectID, http.StatusBadRequest)
		return
	}
	update.BannerID = int32(id)

	err = c.Usecases.UpdateBanner(r.Context(), update)
	if errors.Is(err, BMS.ErrBannerIDNotFound) {
		ProcessError(w, err, http.StatusNotFound)
	} else if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
	}
}

// DeleteBannerHandler handles DELETE request with
// Query params: id;
// Header params: token
func (c *Controller) DeleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		ProcessError(w, ErrIncorrectID, http.StatusBadRequest)
		return
	}

	err = c.Usecases.DeleteBanner(r.Context(), int32(id))
	if errors.Is(err, BMS.ErrBannerIDNotFound) {
		ProcessError(w, err, http.StatusNotFound)
	} else if err != nil {
		ProcessError(w, ErrInternal, http.StatusInternalServerError)
	}
}

const DefaultLimit = 100
const DefaultOffset = 0

func ProcessError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	body := make(map[string]string)
	body["error"] = err.Error()
	bodyJSON, _ := json.Marshal(body)
	w.Write(bodyJSON)
}
