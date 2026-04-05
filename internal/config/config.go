package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Auth     AuthConfig     `mapstructure:"auth"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	AllowOrigins string `mapstructure:"allow_origins"`
}

type AuthConfig struct {
	Type           string        `mapstructure:"type"`
	JWT            JWTConfig     `mapstructure:"jwt"`
	Firebase       FirebaseConfig `mapstructure:"firebase"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

type FirebaseConfig struct {
	CredentialsFile string `mapstructure:"credentials_file"`
	ProjectID       string `mapstructure:"project_id"`
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type CacheConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Type     string        `mapstructure:"type"`
	TTL      time.Duration `mapstructure:"ttl"`
	Cleanup  time.Duration `mapstructure:"cleanup_interval"`
}

type LoggerConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		utils.Warn("No .env file found or error loading .env file: %v", err)
	}

	viper.SetEnvPrefix("LNI")
	viper.AutomaticEnv()

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.allow_origins", "*")

	viper.SetDefault("auth.type", "jwt")
	viper.SetDefault("auth.jwt.secret", "change-this-in-production")
	viper.SetDefault("auth.jwt.access_expiry", "24h")
	viper.SetDefault("auth.jwt.refresh_expiry", "168h")

	viper.SetDefault("firebase.credentials_file", "./firebase-credentials.json")
	viper.SetDefault("firebase.project_id", "")

	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "letsnormalizeit")

	viper.SetDefault("redis.address", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.type", "memory")
	viper.SetDefault("cache.ttl", "15m")
	viper.SetDefault("cache.cleanup_interval", "5m")

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.encoding", "json")
	viper.SetDefault("logger.output_paths", []string{"stdout"})
	viper.SetDefault("logger.error_output_paths", []string{"stderr"})

	configPath := "./configs"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			utils.Warn("Error reading config file: %s", err)
		}
	}

	viper.BindEnv("server.port", "LNI_SERVER_PORT")
	viper.BindEnv("server.allow_origins", "LNI_SERVER_ALLOW_ORIGINS")

	viper.BindEnv("auth.type", "LNI_AUTH_TYPE")
	viper.BindEnv("auth.jwt.secret", "LNI_JWT_SECRET")
	viper.BindEnv("auth.jwt.access_expiry", "LNI_JWT_EXPIRY")
	viper.BindEnv("auth.jwt.refresh_expiry", "LNI_JWT_REFRESH_EXPIRY")

	viper.BindEnv("auth.firebase.credentials_file", "LNI_FIREBASE_CREDENTIALS_FILE")
	viper.BindEnv("auth.firebase.project_id", "LNI_FIREBASE_PROJECT_ID")

	viper.BindEnv("mongodb.uri", "LNI_MONGODB_URI")
	viper.BindEnv("mongodb.database", "LNI_MONGODB_DATABASE")

	viper.BindEnv("redis.address", "LNI_REDIS_ADDRESS")
	viper.BindEnv("redis.password", "LNI_REDIS_PASSWORD")
	viper.BindEnv("redis.db", "LNI_REDIS_DB")

	viper.BindEnv("cache.enabled", "LNI_CACHE_ENABLED")
	viper.BindEnv("cache.type", "LNI_CACHE_TYPE")

	viper.BindEnv("logger.level", "LNI_LOGGER_LEVEL")
	viper.BindEnv("logger.encoding", "LNI_LOGGER_ENCODING")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		utils.Fatal("Error unmarshalling config: %s", err)
	}

	if !filepath.IsAbs(config.Auth.Firebase.CredentialsFile) {
		absPath, _ := filepath.Abs(config.Auth.Firebase.CredentialsFile)
		config.Auth.Firebase.CredentialsFile = absPath
	}

	if err := validateConfig(&config); err != nil {
		utils.Fatal("Configuration validation failed: %s", err)
	}

	return &config
}

func validateConfig(config *Config) error {
	if config.Auth.Type != "jwt" && config.Auth.Type != "firebase" {
		return fmt.Errorf("auth.type must be 'jwt' or 'firebase', got '%s'", config.Auth.Type)
	}

	if config.Auth.Type == "jwt" {
		if config.Auth.JWT.Secret == "" {
			return fmt.Errorf("auth.jwt.secret is required when using JWT auth")
		}
		if config.Auth.JWT.Secret == "change-this-in-production" {
			utils.Warn("Using default JWT secret - change LNI_JWT_SECRET in production")
		}
	}

	if config.Auth.Type == "firebase" {
		if config.Auth.Firebase.ProjectID == "" {
			return fmt.Errorf("auth.firebase.project_id is required when using Firebase auth")
		}
		if config.Auth.Firebase.CredentialsFile == "" {
			return fmt.Errorf("auth.firebase.credentials_file is required when using Firebase auth")
		}
		if _, err := os.Stat(config.Auth.Firebase.CredentialsFile); os.IsNotExist(err) {
			return fmt.Errorf("Firebase credentials file not found: %s", config.Auth.Firebase.CredentialsFile)
		}
	}

	if config.MongoDB.URI == "" {
		return fmt.Errorf("mongodb.uri is required")
	}

	if config.MongoDB.Database == "" {
		return fmt.Errorf("mongodb.database is required")
	}

	return nil
}
