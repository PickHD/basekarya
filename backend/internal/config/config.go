package config

import (
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

func Load() *Config {
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
			Secret:    getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
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
			SuperadminUsername: getEnv("SUPERADMIN_USERNAME", "superadmin"),
			SuperadminPassword: getEnv("SUPERADMIN_PASSWORD", "12390182309812039812093801"),
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
