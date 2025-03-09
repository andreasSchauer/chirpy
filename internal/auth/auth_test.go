package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


func TestJWT(t *testing.T) {
	validSecretKey := "my-secret-test-key"
	correctUserID := uuid.New()

	validToken, err := MakeJWT(correctUserID, validSecretKey, time.Hour)
	if err != nil {
		t.Fatalf("couldn't create tokenString: %v", err)
	}

	expiredToken, err := MakeJWT(correctUserID, validSecretKey, time.Millisecond)
	if err != nil {
		t.Fatalf("couldn't create expired token: %v", err)
	}
	time.Sleep(time.Millisecond * 10)

	tests := []struct {
		name     		string
		tokenSecret   	string
		tokenString		string
		wantUserID		uuid.UUID
		wantErr  		bool
	}{
		{
			name:		  "Correct Key and Token",
			tokenSecret:  validSecretKey,
			tokenString:  validToken,
			wantUserID:   correctUserID,
			wantErr: 	  false,
		},
		{
			name:		  "Wrong Key, Correct Token",
			tokenSecret:  "wrong-key",
			tokenString:  validToken,
			wantUserID:   uuid.Nil,
			wantErr: 	  true,
		},
		{
			name:		  "Invalid token",
			tokenSecret:  validSecretKey,
			tokenString:  "wrong tokenString",
			wantUserID:   uuid.Nil,
			wantErr: 	  true,
		},
		{
			name:		  "Expired token",
			tokenSecret:  validSecretKey,
			tokenString:  expiredToken,
			wantUserID:   uuid.Nil,
			wantErr: 	  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userID, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
			if tc.wantErr {
				if err == nil {
					t.Error("ValidateJWT(): expected error but got none")
				}
				return
			} else {
				if err != nil {
					t.Errorf("ValidateJWT(): unexpected error: %v", err)
					return
				}
				
				if userID != tc.wantUserID {
					t.Errorf("ValidateJWT(): got userID %v, want %v", userID, tc.wantUserID)
				}
			}
		})
	}
}


func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}