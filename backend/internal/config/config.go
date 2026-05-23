package config

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

type Config struct {
	Database              DatabaseConfig
	JWT                   JWTConfig
	Server                ServerConfig
	Logging               LoggingConfig
	Minio                 MinioConfig
	FileUpload            FileUploadConfig
	ExternalServiceConfig ExternalServiceConfig
	CredentialConfig      CredentialConfig
	Redis                 RedisConfig
	Email                 EmailConfig
}

type RedisConfig struct {
	Addr     string
	Password string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn int
}

type ServerConfig struct {
	Port int
	Env  string
}

type LoggingConfig struct {
	Level string
}

type MinioConfig struct {
	Endpoint       string
	AccessKey      string
	SecretKey      string
	BucketName     string
	BucketLocation string
	IsSecure       bool
	PublicDomain   string
}

type FileUploadConfig struct {
	MaxRequestBodySizeMB int
	MaxFileSizeMB        int
}

type ExternalServiceConfig struct {
	NominatimUrl string
}

type CredentialConfig struct {
	SuperadminUsername string
	SuperadminPassword string
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func Load() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	superadminUsername := os.Getenv("SUPERADMIN_USERNAME")
	if superadminUsername == "" {
		panic("SUPERADMIN_USERNAME environment variable is required")
	}
	superadminPassword := os.Getenv("SUPERADMIN_PASSWORD")
	if superadminPassword == "" {
		panic("SUPERADMIN_PASSWORD environment variable is required")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", "root_password"),
			DBName:   getEnv("MYSQL_DATABASE", "hris_db"),
			SSLMode:  getEnv("MYSQL_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:    jwtSecret,
			ExpiresIn: getEnvInt("JWT_EXPIRES_IN_HOUR", 24),
		},
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		Minio: MinioConfig{
			Endpoint:       getEnv("MINIO_ENDPOINT", ""),
			AccessKey:      getEnv("MINIO_ACCESS_KEY", ""),
			SecretKey:      getEnv("MINIO_SECRET_KEY", ""),
			BucketName:     getEnv("MINIO_BUCKET_NAME", ""),
			BucketLocation: getEnv("MINIO_BUCKET_LOCATION", ""),
			IsSecure:       getEnvBool("MINIO_IS_SECURE", false),
			PublicDomain:   getEnv("MINIO_PUBLIC_DOMAIN", ""),
		},
		FileUpload: FileUploadConfig{
			MaxRequestBodySizeMB: getEnvInt("MAX_REQUEST_BODY_SIZE_MB", 50),
			MaxFileSizeMB:        getEnvInt("MAX_FILE_SIZE_MB", 40),
		},
		ExternalServiceConfig: ExternalServiceConfig{
			NominatimUrl: getEnv("NOMINATIM_URL", ""),
		},
		CredentialConfig: CredentialConfig{
			SuperadminUsername: superadminUsername,
			SuperadminPassword: superadminPassword,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", ""),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Email: EmailConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnvInt("SMTP_PORT", 0),
			Username: getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASS", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" {
			return true
		}
		return false
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseInt(value); intValue != 0 {
			return intValue
		}
	}
	return defaultValue
}

func parseInt(s string) int {
	var result int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}

func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable is required", key))
	}
	return value
}
