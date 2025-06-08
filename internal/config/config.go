package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Firebase FirebaseConfig `mapstructure:"firebase"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Redis    RedisConfig    `mapstructure:"redis"`
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

// Load loads the configuration from files and environment variables
func Load() *Config {
	// Set default configuration file path
	configPath := "./configs"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.allow_origins", "*")
	viper.SetDefault("firebase.credentials_file", "./firebase-credentials.json")
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "letsnormalizeit")
	viper.SetDefault("redis.address", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	// Override with environment variables
	viper.SetEnvPrefix("LNI")
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %s", err)
	}

	// Resolve absolute path for Firebase credentials
	if !filepath.IsAbs(config.Firebase.CredentialsFile) {
		absPath, err := filepath.Abs(config.Firebase.CredentialsFile)
		if err == nil {
			config.Firebase.CredentialsFile = absPath
		}
	}

	return &config
}
