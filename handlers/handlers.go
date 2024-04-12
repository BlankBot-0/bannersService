package handlers

import "net/http"

// UserBannerHandler handles GET request with
// Query params: tag_id, feature_id, use_last_version;
// Header params: token
func UserBannerHandler(w http.ResponseWriter, r *http.Request) {

}

// BannersSortedHandler handles GET request with
// Query params: tag_id, feature_id, limit, offset;
// Header params: token
func BannersSortedHandler(w http.ResponseWriter, r *http.Request) {

}

// CreateBannerHandler handles POST request with
// Header params: token
func CreateBannerHandler(w http.ResponseWriter, r *http.Request) {}

// PatchBannerHandler handles PATCH request with
// Query params: id;
// Header params: token
func PatchBannerHandler(w http.ResponseWriter, r *http.Request) {}

// DeleteBannerHandler handles DELETE request with
// Query params: id;
// Header params: token
func DeleteBannerHandler(w http.ResponseWriter, r *http.Request) {}
