package auth

import (
	"testing"
	"time"
)

func TestGenerateTokenPair(t *testing.T) {
	service := NewService("test-secret-key-12345", 1*time.Hour, 7*24*time.Hour)

	tests := []struct {
		name    string
		userID  string
		email   string
		userName string
		isAdmin bool
	}{
		{
			name:    "regular user",
			userID:  "user123",
			email:   "test@example.com",
			userName: "Test User",
			isAdmin: false,
		},
		{
			name:    "admin user",
			userID:  "admin456",
			email:   "admin@example.com",
			userName: "Admin User",
			isAdmin: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := service.GenerateTokenPair(tt.userID, tt.email, tt.userName, tt.isAdmin)
			if err != nil {
				t.Fatalf("GenerateTokenPair() error = %v", err)
			}

			if tokens.AccessToken == "" {
				t.Error("AccessToken should not be empty")
			}

			if tokens.RefreshToken == "" {
				t.Error("RefreshToken should not be empty")
			}

			if tokens.ExpiresIn <= 0 {
				t.Error("ExpiresIn should be positive")
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	service := NewService("test-secret-key-12345", 1*time.Hour, 7*24*time.Hour)

	t.Run("valid token", func(t *testing.T) {
		userID := "user123"
		email := "test@example.com"
		name := "Test User"
		isAdmin := true

		tokens, _ := service.GenerateTokenPair(userID, email, name, isAdmin)

		claims, err := service.VerifyToken(tokens.AccessToken)
		if err != nil {
			t.Fatalf("VerifyToken() error = %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("UserID = %s, want %s", claims.UserID, userID)
		}

		if claims.Email != email {
			t.Errorf("Email = %s, want %s", claims.Email, email)
		}

		if claims.Name != name {
			t.Errorf("Name = %s, want %s", claims.Name, name)
		}

		if claims.IsAdmin != isAdmin {
			t.Errorf("IsAdmin = %v, want %v", claims.IsAdmin, isAdmin)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := service.VerifyToken("invalid-token")
		if err == nil {
			t.Error("VerifyToken() expected error for invalid token")
		}
	})

	t.Run("token with wrong secret", func(t *testing.T) {
		otherService := NewService("different-secret", 1*time.Hour, 7*24*time.Hour)
		tokens, _ := otherService.GenerateTokenPair("user", "email", "name", false)

		_, err := service.VerifyToken(tokens.AccessToken)
		if err == nil {
			t.Error("VerifyToken() expected error for token with wrong secret")
		}
	})
}

func TestRefreshToken(t *testing.T) {
	service := NewService("test-secret-key-12345", 1*time.Hour, 7*24*time.Hour)

	userID := "user123"
	email := "test@example.com"
	name := "Test User"
	isAdmin := false

	originalTokens, _ := service.GenerateTokenPair(userID, email, name, isAdmin)

	newTokens, err := service.RefreshToken(originalTokens.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Error("New AccessToken should not be empty")
	}

	claims, _ := service.VerifyToken(newTokens.AccessToken)
	if claims.UserID != userID {
		t.Errorf("UserID = %s, want %s", claims.UserID, userID)
	}
}

func TestPasswordHashing(t *testing.T) {
	service := NewService("test-secret-key", 1*time.Hour, 7*24*time.Hour)

	t.Run("hash and check password", func(t *testing.T) {
		password := "securepassword123"

		hash, err := service.HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword() error = %v", err)
		}

		if hash == password {
			t.Error("Hash should not equal plain password")
		}

		if !service.CheckPassword(password, hash) {
			t.Error("CheckPassword() should return true for correct password")
		}

		if service.CheckPassword("wrongpassword", hash) {
			t.Error("CheckPassword() should return false for wrong password")
		}
	})
}

func TestGenerateUserID(t *testing.T) {
	id1 := GenerateUserID()
	id2 := GenerateUserID()

	if id1 == "" {
		t.Error("GenerateUserID() should not return empty string")
	}

	if id1 == id2 {
		t.Error("GenerateUserID() should generate unique IDs")
	}

	if len(id1) != 24 {
		t.Errorf("GenerateUserID() should return 24 character hex string, got %d", len(id1))
	}
}
