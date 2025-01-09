package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/SeeALie/user/internal/auth"
	"github.com/jaydee029/SeeALie/user/internal/database"
	validate "github.com/jaydee029/SeeALie/user/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserRes struct {
	Username   string    `json:"username"`
	Created_at time.Time `json:"created_at"`
	Email      string    `json:"email"`
}
type res_login struct {
	Email         string `json:"email"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}
type UserInput struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (cfg *apiconfig) Signup(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := UserInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	err = validate.ValidateEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	email_if_exist, err := cfg.DB.Is_email(context.Background(), params.Email)

	if email_if_exist {
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.ValidateUsername(params.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	username_if_exists, err := cfg.DB.Is_username(r.Context(), params.Username)
	if username_if_exists {
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.ValidatePassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	encrypted, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	uuids := uuid.New().String()
	var pgUUID pgtype.UUID

	err = pgUUID.Scan(uuids)
	if err != nil {
		fmt.Println("Error setting UUID:", err)
		return
	}

	var pgtime pgtype.Timestamp
	err = pgtime.Scan(time.Now().UTC())
	if err != nil {
		fmt.Println("Error setting timestamp:", err)
		return
	}

	user, err := cfg.DB.Createuser(r.Context(), database.CreateuserParams{
		Email:     params.Email,
		Username:  params.Username,
		Passwd:    encrypted,
		ID:        pgUUID,
		CreatedAt: pgtime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, UserRes{
		Created_at: user.CreatedAt.Time,
		Username:   user.Username,
		Email:      user.Email,
	})
}

func (cfg *apiconfig) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := UserInput{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	var identifier string
	if err := validate.ValidateEmail(params.Email); err != nil {
		if err := validate.ValidateUsername(params.Username); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid email or username")
			return
		}
		identifier = params.Username
	}
	identifier = params.Email

	if identifier == params.Email {

	} else if identifier == params.Username {
	}

	user, err := cfg.DB.Find_user_email(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password doesn't match")
		return
	}

	Userid, _ := uuid.FromBytes(user.ID.Bytes[:])

	Token, err := auth.Tokenize(Userid, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Refresh_token, err := auth.RefreshToken(Userid, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, res_login{
		Email:         params.Email,
		Token:         Token,
		Refresh_token: Refresh_token,
	})

}
