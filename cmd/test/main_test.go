package test

import (
	"banners/internal/handlers"
	"banners/internal/usecase"
	bytes "bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type tokenPayload struct {
	Token string `json:"token"`
}

func TestAuth(t *testing.T) {
	if pool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432"); err != nil {
		panic(err)
	} else if _, err = pool.Exec(context.Background(), "truncate banners, banners_tag, banners_info"); err != nil {
		panic(err)
	}

	url := "http://localhost:8080"
	client := http.DefaultClient

	var featureID int32 = 102

	var token tokenPayload
	{
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/admin_token?username=%s&password=%s", url, "admin", "password"), nil)
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		err = json.NewDecoder(resp.Body).Decode(&token)
		assert.NoError(t, err)
		assert.NotEmpty(t, token.Token)
	}

	dto := usecase.BannerDTO{
		FeatureID: featureID,
		Tags:      []int32{101},
		Contents:  `{"kek": 1}`,
		IsActive:  true,
	}
	{
		body, err := json.Marshal(dto)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, url+"/banner", bytes.NewBuffer(body))
		assert.NoError(t, err)

		req.Header.Set("Authorization", "Bearer "+token.Token)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	{
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_banner?feature_id=%d&tag_id=%d", url, featureID, 101), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token.Token)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var content string
		err = json.NewDecoder(resp.Body).Decode(&content)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Equal(t, dto.Contents, content)
	}
}

func TestProcessError(t *testing.T) {
	w := httptest.NewRecorder()
	err := handlers.ErrIncorrectBannerContent
	code := http.StatusBadRequest

	handlers.ProcessError(w, err, code)
	got := w.Body.String()

	want := "{\"error\":\"incorrect banner content\"}"
	if got != want {
		t.Errorf("incorrect response body: got %q, want %q", got, want)
	}
	if w.Result().StatusCode != code {
		t.Errorf("incorrect response status: got %d, want %d", w.Result().StatusCode, code)
	}
}
