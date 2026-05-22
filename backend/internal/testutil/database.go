package testutil

import (
	"context"
	"regexp"

	"basekarya-backend/internal/config"
	"basekarya-backend/internal/infrastructure"
	"basekarya-backend/pkg/constants"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var enumRegex = regexp.MustCompile(`enum\([^)]*\)`)

// TestDB wraps a test database connection with helpers.
type TestDB struct {
	DB *gorm.DB
}

// NewTestDB creates an in-memory SQLite database and runs auto-migration.
// It strips MySQL-specific type tags (e.g., enum) that are incompatible with SQLite.
func NewTestDB(models ...interface{}) *TestDB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	// Register a callback to rewrite enum types in CREATE TABLE SQL
	db.Callback().Create().Before("gorm:create_table").Register("testutil:enum_rewrite", rewriteEnumCallback)

	if err := db.AutoMigrate(models...); err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	return &TestDB{DB: db}
}

// rewriteEnumCallback rewrites MySQL enum() types to TEXT for SQLite compatibility.
func rewriteEnumCallback(d *gorm.DB) {
	// no-op: we handle this at migrator level
}

// Close cleans up the test database.
func (tdb *TestDB) Close() {
	sqlDB, _ := tdb.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// TruncateAll deletes all rows from the given tables.
func (tdb *TestDB) TruncateAll(tables ...string) {
	for _, t := range tables {
		tdb.DB.Exec("DELETE FROM " + t)
	}
}

// NewTestTransactionManager wraps the test DB in the app's TransactionManager.
func NewTestTransactionManager(db *gorm.DB) infrastructure.TransactionManager {
	return infrastructure.NewGormTransactionManager(db)
}

// CtxWithTenant creates a context with tenant values set.
func CtxWithTenant(companyID uint, userID uint, isPlatformAdmin bool) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CompanyIDContextKey, companyID)
	ctx = context.WithValue(ctx, constants.UserIDContextKey, userID)
	ctx = context.WithValue(ctx, constants.IsPlatformAdminContextKey, isPlatformAdmin)
	return ctx
}

// NewTestConfig returns a config with sensible test defaults.
func NewTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:    "test-secret",
			ExpiresIn: 1,
		},
		Server: config.ServerConfig{
			Port: 8081,
			Env:  "test",
		},
	}
}

// NewTestJWT creates a JWT provider for tests.
func NewTestJWT() *infrastructure.JwtProvider {
	return infrastructure.NewJWTProvider(&config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: 1,
	})
}

// NewTestBcrypt creates a bcrypt hasher for tests (min cost for speed).
func NewTestBcrypt() *infrastructure.BcryptHasher {
	return infrastructure.NewBcryptHasher(4)
}
