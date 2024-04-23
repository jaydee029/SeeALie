package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/jaydee029/SeeALie/authentication/internal"
	"github.com/jaydee029/SeeALie/authentication/internal/database"
)

type Refresh_res struct {
	Refresh_token string    `json:"refresh_token,omitempty"`
	Auth_token    string    `json:"auth_token,omitempty"`
	Revoked_at    time.Time `json:"revoked_at,omitempty"`
}

func (cfg *apiconfig) Refresh(w http.ResponseWriter, r *http.Request) {
	Token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	is_refresh, _ := auth.VerifyRefresh(Token, cfg.jwtsecret)

	if !is_refresh {
		respondWithError(w, http.StatusBadRequest, "invalid refresh token")
		return
	}

	user_id, err := auth.ValidateToken(Token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := uuid.Parse(user_id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	auth_token, err := auth.Tokenize(id, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, Refresh_res{
		Auth_token: auth_token,
	})

}
func (cfg *apiconfig) Revoke(w http.ResponseWriter, r *http.Request) {

	Token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	is_refresh, _ := auth.VerifyRefresh(Token, cfg.jwtsecret)

	if !is_refresh {
		respondWithError(w, http.StatusBadRequest, "invalid refresh token")
		return
	}

	revoked, err := cfg.DB.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     Token,
		RevokedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, Refresh_res{
		Refresh_token: revoked.Token,
		Revoked_at:    revoked.RevokedAt,
	})

}
