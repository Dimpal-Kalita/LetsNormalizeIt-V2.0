package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	os.Setenv("LNI_AUTH_TYPE", "jwt")
	os.Setenv("LNI_JWT_SECRET", "test-secret")
	os.Setenv("LNI_JWT_EXPIRY", "2h")
	os.Setenv("LNI_MONGODB_URI", "mongodb://test:27017")
	os.Setenv("LNI_MONGODB_DATABASE", "testdb")
	os.Setenv("LNI_CACHE_ENABLED", "false")

	cfg := Load()

	if cfg.Auth.Type != "jwt" {
		t.Errorf("Auth.Type = %s, want jwt", cfg.Auth.Type)
	}

	if cfg.Auth.JWT.Secret != "test-secret" {
		t.Errorf("Auth.JWT.Secret = %s, want test-secret", cfg.Auth.JWT.Secret)
	}

	if cfg.Auth.JWT.AccessExpiry != 2*time.Hour {
		t.Errorf("Auth.JWT.AccessExpiry = %v, want 2h", cfg.Auth.JWT.AccessExpiry)
	}

	if cfg.MongoDB.URI != "mongodb://test:27017" {
		t.Errorf("MongoDB.URI = %s, want mongodb://test:27017", cfg.MongoDB.URI)
	}

	if cfg.MongoDB.Database != "testdb" {
		t.Errorf("MongoDB.Database = %s, want testdb", cfg.MongoDB.Database)
	}

	if cfg.Cache.Enabled != false {
		t.Errorf("Cache.Enabled = %v, want false", cfg.Cache.Enabled)
	}
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid JWT config", func(t *testing.T) {
		cfg := &Config{
			Auth: AuthConfig{
				Type: "jwt",
				JWT: JWTConfig{
					Secret:        "test-secret",
					AccessExpiry:  24 * time.Hour,
					RefreshExpiry: 168 * time.Hour,
				},
			},
			MongoDB: MongoDBConfig{
				URI:      "mongodb://localhost:27017",
				Database: "testdb",
			},
		}

		if err := validateConfig(cfg); err != nil {
			t.Errorf("validateConfig() error = %v", err)
		}
	})

	t.Run("missing JWT secret", func(t *testing.T) {
		cfg := &Config{
			Auth: AuthConfig{
				Type: "jwt",
				JWT: JWTConfig{
					Secret: "",
				},
			},
			MongoDB: MongoDBConfig{
				URI:      "mongodb://localhost:27017",
				Database: "testdb",
			},
		}

		if err := validateConfig(cfg); err == nil {
			t.Error("validateConfig() expected error for missing JWT secret")
		}
	})

	t.Run("missing MongoDB URI", func(t *testing.T) {
		cfg := &Config{
			Auth: AuthConfig{
				Type: "jwt",
				JWT: JWTConfig{
					Secret: "test-secret",
				},
			},
			MongoDB: MongoDBConfig{
				URI: "",
			},
		}

		if err := validateConfig(cfg); err == nil {
			t.Error("validateConfig() expected error for missing MongoDB URI")
		}
	})

	t.Run("invalid auth type", func(t *testing.T) {
		cfg := &Config{
			Auth: AuthConfig{
				Type: "invalid",
			},
		}

		if err := validateConfig(cfg); err == nil {
			t.Error("validateConfig() expected error for invalid auth type")
		}
	})
}
