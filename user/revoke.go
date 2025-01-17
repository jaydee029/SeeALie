package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/SeeALie/user/internal/auth"
	"github.com/jaydee029/SeeALie/user/internal/database"
	"github.com/jaydee029/SeeALie/user/utils"
)

type refreshRes struct {
	Refresh_token string           `json:"refresh_token,omitempty"`
	Auth_token    string           `json:"auth_token,omitempty"`
	Auth_Expiry   pgtype.Timestamp `json:"auth_expiry,omitempty"`
	Revoked_at    pgtype.Timestamp `json:"revoked_at,omitempty"`
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

	useridstr, err := auth.ValidateToken(Token, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userid, err := utils.GenpgtypeUUID(useridstr)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	auth_token, expiresat, err := auth.Tokenize(userid, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expiresatpgtype, err := utils.GenpgtypeTimestamp(expiresat)
	if err != nil {
		log.Println("Error setting timestamp:", err)
	}

	respondWithJson(w, http.StatusAccepted, refreshRes{
		Auth_token:  auth_token,
		Auth_Expiry: expiresatpgtype,
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

	revokedatpgtype, err := utils.GenpgtypeTimestamp(time.Now().UTC())
	if err != nil {
		log.Println("Error setting timestamp:", err)
	}

	revoked, err := cfg.DB.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token:     Token,
		RevokedAt: revokedatpgtype,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, refreshRes{
		Refresh_token: revoked.Token,
		Revoked_at:    revoked.RevokedAt,
	})

}
