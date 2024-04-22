package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/jaydee029/SeeALie/authentication/internal"
	"github.com/jaydee029/SeeALie/authentication/internal/database"
	validate "github.com/jaydee029/SeeALie/authentication/validator"
	"golang.org/x/crypto/bcrypt"
)

type Input struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Passwd   string `json:"passwd"`
}

type login_input struct {
	Login_id string `json:"login_id"` //email or username
	Passwd   string `json:"passwd"`
}

type Res struct {
	Username      string    `json:"username,omitempty"`
	Refresh_Token string    `json:"refresh_token,omitempty"`
	Auth_Token    string    `json:"auth_token,omitempty"`
	Created_at    time.Time `json:"created_at,omitempty"`
}

func (cfg *apiconfig) signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Input{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reading the input")
		return
	}

	is_email, err := validate.ValidateEmail(params.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !is_email {
		respondWithError(w, http.StatusBadRequest, "Enter a valid email")
		return
	}

	is_valid_name, err := validate.ValidateUsername(params.Username)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !is_valid_name {
		respondWithError(w, http.StatusBadRequest, "Enter a valid username")
		return
	}

	is_valid_passwd, err := validate.ValidatePassword(params.Passwd)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !is_valid_passwd {
		respondWithError(w, http.StatusBadRequest, "Enter a valid password")
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

	respondWithJson(w, http.StatusCreated, Res{
		Username:   user.Username,
		Created_at: user.CreatedAt,
	})

}

func (cfg *apiconfig) login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := login_input{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	var user database.User

	is_email, _ := validate.ValidateEmail(params.Login_id)

	if is_email {
		user, err = cfg.DB.Find_user_email(r.Context(), params.Login_id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	} else {
		is_valid_name, _ := validate.ValidateUsername(params.Login_id)
		if is_valid_name {
			user, err = cfg.DB.Find_user_name(r.Context(), params.Login_id)
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}
	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(params.Passwd))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password doesn't match")
		return
	}

	auth_token, err := auth.Tokenize(user.ID, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refresh_token, err := auth.RefreshToken(user.ID, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, Res{
		Username:      user.Username,
		Refresh_Token: refresh_token,
		Auth_Token:    auth_token,
	})

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
