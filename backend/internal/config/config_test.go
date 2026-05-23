package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	envKeys := []string{
		"MYSQL_HOST", "MYSQL_PORT", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DATABASE", "MYSQL_SSLMODE",
		"JWT_SECRET", "JWT_EXPIRES_IN_HOUR",
		"SERVER_PORT", "SERVER_ENV",
		"LOG_LEVEL",
		"MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "MINIO_BUCKET_NAME", "MINIO_BUCKET_LOCATION", "MINIO_IS_SECURE", "MINIO_PUBLIC_DOMAIN",
		"MAX_REQUEST_BODY_SIZE_MB", "MAX_FILE_SIZE_MB",
		"NOMINATIM_URL",
		"SUPERADMIN_USERNAME", "SUPERADMIN_PASSWORD",
		"REDIS_ADDR", "REDIS_PASSWORD",
		"SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASS", "SMTP_FROM",
	}
	for _, k := range envKeys {
		os.Unsetenv(k)
	}

	cfg := Load()

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected Database.Host localhost, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != "3306" {
		t.Errorf("expected Database.Port 3306, got %s", cfg.Database.Port)
	}
	if cfg.Database.User != "root" {
		t.Errorf("expected Database.User root, got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "root_password" {
		t.Errorf("expected Database.Password root_password, got %s", cfg.Database.Password)
	}
	if cfg.Database.DBName != "hris_db" {
		t.Errorf("expected Database.DBName hris_db, got %s", cfg.Database.DBName)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected Database.SSLMode disable, got %s", cfg.Database.SSLMode)
	}
	if cfg.JWT.Secret != "your-super-secret-jwt-key" {
		t.Errorf("expected JWT.Secret your-super-secret-jwt-key, got %s", cfg.JWT.Secret)
	}
	if cfg.JWT.ExpiresIn != 24 {
		t.Errorf("expected JWT.ExpiresIn 24, got %d", cfg.JWT.ExpiresIn)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.Env != "development" {
		t.Errorf("expected Server.Env development, got %s", cfg.Server.Env)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("expected Logging.Level debug, got %s", cfg.Logging.Level)
	}
	if cfg.Minio.IsSecure != false {
		t.Errorf("expected Minio.IsSecure false, got %v", cfg.Minio.IsSecure)
	}
	if cfg.FileUpload.MaxRequestBodySizeMB != 50 {
		t.Errorf("expected FileUpload.MaxRequestBodySizeMB 50, got %d", cfg.FileUpload.MaxRequestBodySizeMB)
	}
	if cfg.FileUpload.MaxFileSizeMB != 40 {
		t.Errorf("expected FileUpload.MaxFileSizeMB 40, got %d", cfg.FileUpload.MaxFileSizeMB)
	}
	if cfg.Email.Port != 0 {
		t.Errorf("expected Email.Port 0, got %d", cfg.Email.Port)
	}
}

func TestLoad_EnvVars(t *testing.T) {
	os.Setenv("MYSQL_HOST", "dbhost")
	os.Setenv("MYSQL_PORT", "5432")
	os.Setenv("JWT_SECRET", "mysecret")
	os.Setenv("JWT_EXPIRES_IN_HOUR", "48")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_ENV", "production")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("MINIO_IS_SECURE", "true")
	os.Setenv("SMTP_PORT", "587")
	t.Cleanup(func() {
		os.Unsetenv("MYSQL_HOST")
		os.Unsetenv("MYSQL_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRES_IN_HOUR")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_ENV")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("MINIO_IS_SECURE")
		os.Unsetenv("SMTP_PORT")
	})

	cfg := Load()

	if cfg.Database.Host != "dbhost" {
		t.Errorf("expected Database.Host dbhost, got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != "5432" {
		t.Errorf("expected Database.Port 5432, got %s", cfg.Database.Port)
	}
	if cfg.JWT.Secret != "mysecret" {
		t.Errorf("expected JWT.Secret mysecret, got %s", cfg.JWT.Secret)
	}
	if cfg.JWT.ExpiresIn != 48 {
		t.Errorf("expected JWT.ExpiresIn 48, got %d", cfg.JWT.ExpiresIn)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected Server.Port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.Env != "production" {
		t.Errorf("expected Server.Env production, got %s", cfg.Server.Env)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("expected Logging.Level info, got %s", cfg.Logging.Level)
	}
	if cfg.Minio.IsSecure != true {
		t.Errorf("expected Minio.IsSecure true, got %v", cfg.Minio.IsSecure)
	}
	if cfg.Email.Port != 587 {
		t.Errorf("expected Email.Port 587, got %d", cfg.Email.Port)
	}
}

func TestGetEnv(t *testing.T) {
	os.Unsetenv("TEST_GETENV_KEY")
	result := getEnv("TEST_GETENV_KEY", "fallback")
	if result != "fallback" {
		t.Errorf("expected fallback, got %s", result)
	}

	os.Setenv("TEST_GETENV_KEY", "myval")
	t.Cleanup(func() { os.Unsetenv("TEST_GETENV_KEY") })

	result = getEnv("TEST_GETENV_KEY", "fallback")
	if result != "myval" {
		t.Errorf("expected myval, got %s", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Unsetenv("TEST_GETENVINT_KEY")
	result := getEnvInt("TEST_GETENVINT_KEY", 42)
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}

	os.Setenv("TEST_GETENVINT_KEY", "100")
	t.Cleanup(func() { os.Unsetenv("TEST_GETENVINT_KEY") })

	result = getEnvInt("TEST_GETENVINT_KEY", 42)
	if result != 100 {
		t.Errorf("expected 100, got %d", result)
	}
}

func TestGetEnvBool(t *testing.T) {
	os.Unsetenv("TEST_GETENVBOOL_KEY")
	if got := getEnvBool("TEST_GETENVBOOL_KEY", false); got != false {
		t.Errorf("expected false, got %v", got)
	}
	if got := getEnvBool("TEST_GETENVBOOL_KEY", true); got != true {
		t.Errorf("expected true, got %v", got)
	}

	os.Setenv("TEST_GETENVBOOL_KEY", "true")
	t.Cleanup(func() { os.Unsetenv("TEST_GETENVBOOL_KEY") })
	if got := getEnvBool("TEST_GETENVBOOL_KEY", false); got != true {
		t.Errorf("expected true for 'true', got %v", got)
	}

	os.Setenv("TEST_GETENVBOOL_KEY", "1")
	if got := getEnvBool("TEST_GETENVBOOL_KEY", false); got != true {
		t.Errorf("expected true for '1', got %v", got)
	}

	os.Setenv("TEST_GETENVBOOL_KEY", "false")
	if got := getEnvBool("TEST_GETENVBOOL_KEY", true); got != false {
		t.Errorf("expected false for 'false', got %v", got)
	}

	os.Setenv("TEST_GETENVBOOL_KEY", "0")
	if got := getEnvBool("TEST_GETENVBOOL_KEY", true); got != false {
		t.Errorf("expected false for '0', got %v", got)
	}

	os.Setenv("TEST_GETENVBOOL_KEY", "random")
	if got := getEnvBool("TEST_GETENVBOOL_KEY", true); got != false {
		t.Errorf("expected false for 'random', got %v", got)
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"9999", 9999},
		{"12a3", 0},
		{"abc", 0},
		{"", 0},
		{"-5", 0},
		{"42", 42},
	}
	for _, tt := range tests {
		got := parseInt(tt.input)
		if got != tt.expected {
			t.Errorf("parseInt(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}
