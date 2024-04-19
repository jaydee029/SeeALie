package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/jaydee029/SeeALie/authentication/internal"
	"github.com/jaydee029/SeeALie/authentication/internal/database"
)

type Input struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

type Res struct {
	Username      string `json:"username"`
	Refresh_Token string `json:"refresh_token"`
}

func (cfg *apiconfig) signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Input{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading the input")
		return
	}

	if_email, err := cfg.DB.If_email(r.Context(), params.Email)

	if if_email {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if_username, err := cfg.DB.If_username(r.Context(), params.Username)

	if if_username {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	encrypted, err := auth.Hashpassword(params.Passwd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id := uuid.New()

	user, err := cfg.DB.Createuser(r.Context(), database.CreateuserParams{
		ID:        id,
		Username:  params.Username,
		Email:     params.Email,
		Passwd:    encrypted,
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refresh_token, err := auth.RefreshToken(id, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, Res{
		Username:      user.Username,
		Refresh_Token: refresh_token,
	})

}

func (cfg *apiconfig) login(w http.ResponseWriter, r *http.Request) {

	token, err := auth.BearerHeader(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	is_refresh, err := auth.VerifyRefresh(token, cfg.jwtsecret)

	if !is_refresh {
		respondWithError(w, http.StatusBadRequest, "not a valid refresh token")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

func respondWithError(w http.ResponseWriter, code int, res string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", res)
	}
	type errresponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errresponse{
		Error: res,
	})
}

func respondWithJson(w http.ResponseWriter, code int, res interface{}) {
	w.Header().Set("content-type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
