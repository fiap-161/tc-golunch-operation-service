package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"normal password", "mypassword", false},
		{"empty password", "", false},
		{"long password", "thisisareallylongpassword1234567890!@#$%^&*()", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if hash == "" && !tt.wantErr {
				t.Error("HashPassword() returned empty hash without error")
			}

			if !CheckPasswordHash(tt.password, hash) {
				t.Errorf("CheckPasswordHash() = false, want true")
			}

			if CheckPasswordHash("wrongpassword", hash) {
				t.Errorf("CheckPasswordHash() with wrong password = true, want false")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	validHash, err := HashPassword("correctpassword")
	if err != nil {
		t.Fatalf("Failed to hash password for tests: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"correct password", "correctpassword", validHash, true},
		{"incorrect password", "wrongpassword", validHash, false},
		{"empty password", "", validHash, false},
		{"empty hash", "correctpassword", "", false},
		{"both empty", "", "", false},
		{"invalid hash format", "any", "notahash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordHash(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
