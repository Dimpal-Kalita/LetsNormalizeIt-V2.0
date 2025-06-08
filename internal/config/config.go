package config

import (
	"os"
	"path/filepath"

	"github.com/dksensei/letsnormalizeit/internal/utils"
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

	// Logger defaults
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.encoding", "json")
	viper.SetDefault("logger.output_paths", []string{"stdout"})
	viper.SetDefault("logger.error_output_paths", []string{"stderr"})

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			utils.Warn("Config file not found, using defaults and environment variables")
		} else {
			utils.Fatal("Error reading config file: %s", err)
		}
	}

	// Override with environment variables
	viper.SetEnvPrefix("LNI")
	viper.AutomaticEnv()

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

	return &config
}
