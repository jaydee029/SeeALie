package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/SeeALie/user/internal/auth"
	"github.com/jaydee029/SeeALie/user/internal/database"
	"github.com/jaydee029/SeeALie/user/utils"
	validate "github.com/jaydee029/SeeALie/user/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserRes struct {
	Username   string    `json:"username"`
	Created_at time.Time `json:"created_at"`
	Email      string    `json:"email"`
}
type LoginRes struct {
	Email         string           `json:"email"`
	Token         string           `json:"token"`
	Refresh_token string           `json:"refresh_token"`
	ExpiresAt     pgtype.Timestamp `json:"expiresat"`
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

	email_if_exist, err := cfg.DB.IfEmail(context.Background(), params.Email)

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

	username_if_exists, err := cfg.DB.IfUsername(r.Context(), params.Username)
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

	userid, err := utils.GenpgtypeUUID(uuid.New().String())
	if err != nil {
		log.Println("Error setting UUID:", err)
		return
	}

	createdat, err := utils.GenpgtypeTimestamp(time.Now().UTC())
	if err != nil {
		fmt.Println("Error setting timestamp:", err)
		return
	}

	user, err := cfg.DB.Createuser(r.Context(), database.CreateuserParams{
		Email:     params.Email,
		Username:  params.Username,
		Passwd:    encrypted,
		ID:        userid,
		CreatedAt: createdat,
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

	var user database.User
	if identifier == params.Email {
		user, err = cfg.DB.FindUserByEmail(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	} else if identifier == params.Username {
		user, err = cfg.DB.FindUserByUsername(r.Context(), params.Username)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	ifsession, err := cfg.DB.FindSessionByid(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ifsession {
		respondWithError(w, http.StatusUnauthorized, "User already logged in")
	}
	err = bcrypt.CompareHashAndPassword(user.Passwd, []byte(params.Password))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password doesn't match")
		return
	}

	Token, expiresat, err := auth.Tokenize(user.ID, cfg.jwtsecret)

	expiresatpgtype, err := utils.GenpgtypeTimestamp(expiresat)
	if err != nil {
		log.Println("Error setting timestamp:", err)
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	Refresh_token, err := auth.RefreshToken(user.ID, cfg.jwtsecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	session, err := cfg.DB.InsertSession(r.Context(), database.InsertSessionParams{
		SessionID: user.ID,
		UserID:    user.ID,
		ExpiresAt: expiresatpgtype,
		Jwt:       Token,
	})

	respondWithJson(w, http.StatusOK, LoginRes{
		Email:         params.Email,
		Token:         Token,
		ExpiresAt:     session.ExpiresAt,
		Refresh_token: Refresh_token,
	})

}
