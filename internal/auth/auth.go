package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)


func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't generate hash from password: %v", err)
	}

	return string(hash), nil
}


func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	now := time.Now().UTC()
	expirationTime := now.Add(expiresIn)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: 	string(TokenTypeAccess),
		IssuedAt: 	jwt.NewNumericDate(now),
		ExpiresAt: 	jwt.NewNumericDate(expirationTime),
		Subject: 	userID.String(),
	})

	return token.SignedString(signingKey)
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %v", err)
	}

	return uuid.UUID(userID), nil
}


func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization"); 

	if authHeader == "" {
		return "", errors.New("no authorization included in request")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("no bearer token found")
	}

	bearerToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	
	return bearerToken, nil
}