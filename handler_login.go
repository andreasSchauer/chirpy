package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andreasSchauer/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 			string 	`json:"password"`
		Email 				string 	`json:"email"`
		ExpiresInSeconds	int 	`json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context() ,params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No user with this email exists", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password", err)
		return
	}

	expirationTime := params.ExpiresInSeconds
	defaultExpirationTime := 3600

	if expirationTime == 0 || expirationTime > defaultExpirationTime {
		expirationTime = defaultExpirationTime
	}

	expirationDuration := time.Duration(expirationTime) * time.Second

	accessToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, expirationDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Email:     user.Email,
		},
		Token:	   accessToken,
	})
}