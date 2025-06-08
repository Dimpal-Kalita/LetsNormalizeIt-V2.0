package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Firebase FirebaseConfig `mapstructure:"firebase"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string `mapstructure:"port"`
	AllowOrigins string `mapstructure:"allow_origins"`
}

// FirebaseConfig holds Firebase-specific configuration
type FirebaseConfig struct {
	CredentialsFile string `mapstructure:"credentials_file"`
	ProjectID       string `mapstructure:"project_id"`
}

// MongoDBConfig holds MongoDB-specific configuration
type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LoggerConfig holds logger-specific configuration
type LoggerConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

// Load loads the configuration from files and environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		utils.Warn("No .env file found or error loading .env file: %v", err)
	}

	// Set environment variable prefix and enable automatic env reading
	viper.SetEnvPrefix("LNI")
	viper.AutomaticEnv()

	// Set default values (will be overridden by env vars if present)
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.allow_origins", "*")
	viper.SetDefault("firebase.credentials_file", "./firebase-credentials.json")
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "letsnormalizeit")
	viper.SetDefault("redis.address", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.encoding", "json")
	viper.SetDefault("logger.output_paths", []string{"stdout"})
	viper.SetDefault("logger.error_output_paths", []string{"stderr"})

	// Try to read config file as fallback (optional)
	configPath := "./configs"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Read the config file (optional - env vars take precedence)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			utils.Info("Config file not found, using environment variables and defaults")
		} else {
			utils.Warn("Error reading config file: %s, using environment variables and defaults", err)
		}
	} else {
		utils.Info("Config file loaded from: %s", viper.ConfigFileUsed())
	}

	// Explicitly bind environment variables to ensure they override config file values
	viper.BindEnv("server.port", "LNI_SERVER_PORT")
	viper.BindEnv("server.allow_origins", "LNI_SERVER_ALLOW_ORIGINS")
	viper.BindEnv("firebase.credentials_file", "LNI_FIREBASE_CREDENTIALS_FILE")
	viper.BindEnv("firebase.project_id", "LNI_FIREBASE_PROJECT_ID")
	viper.BindEnv("mongodb.uri", "LNI_MONGODB_URI")
	viper.BindEnv("mongodb.database", "LNI_MONGODB_DATABASE")
	viper.BindEnv("redis.address", "LNI_REDIS_ADDRESS")
	viper.BindEnv("redis.password", "LNI_REDIS_PASSWORD")
	viper.BindEnv("redis.db", "LNI_REDIS_DB")
	viper.BindEnv("logger.level", "LNI_LOGGER_LEVEL")
	viper.BindEnv("logger.encoding", "LNI_LOGGER_ENCODING")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		utils.Fatal("Error unmarshalling config: %s", err)
	}

	// Resolve absolute path for Firebase credentials
	if !filepath.IsAbs(config.Firebase.CredentialsFile) {
		absPath, err := filepath.Abs(config.Firebase.CredentialsFile)
		if err == nil {
			config.Firebase.CredentialsFile = absPath
		}
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		utils.Fatal("Configuration validation failed: %s", err)
	}

	return &config
}

func validateConfig(config *Config) error {
	if config.Firebase.ProjectID == "" {
		return fmt.Errorf("Firebase project ID is required")
	}

	if config.Firebase.CredentialsFile == "" {
		return fmt.Errorf("Firebase credentials file path is required")
	}

	if config.MongoDB.URI == "" {
		return fmt.Errorf("MongoDB URI is required")
	}

	if config.MongoDB.Database == "" {
		return fmt.Errorf("MongoDB database name is required")
	}

	// if config.Redis.Address == "" {
	// 	return fmt.Errorf("Redis address is required")
	// }

	// Check if Firebase credentials file exists
	if _, err := os.Stat(config.Firebase.CredentialsFile); os.IsNotExist(err) {
		return fmt.Errorf("Firebase credentials file not found: %s", config.Firebase.CredentialsFile)
	}

	return nil
}

