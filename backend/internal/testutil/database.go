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
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

var enumRegex = regexp.MustCompile(`enum\([^)]*\)`)

// TestDB wraps a test database connection with helpers.
type TestDB struct {
	DB *gorm.DB
}

// NewTestDB creates an in-memory SQLite database and runs auto-migration.
// It strips MySQL-specific type tags (e.g., enum) that are incompatible with SQLite.
func NewTestDB(models ...interface{}) *TestDB {
	db, err := gorm.Open(SQLiteEnumDialector{sqlite.Open(":memory:")}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	if err := db.AutoMigrate(models...); err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	return &TestDB{DB: db}
}

// SQLiteEnumDialector wraps the SQLite dialector to rewrite MySQL enum types as TEXT.
type SQLiteEnumDialector struct {
	gorm.Dialector
}

func (d SQLiteEnumDialector) DataTypeOf(field *schema.Field) string {
	if enumRegex.MatchString(string(field.DataType)) {
		return "text"
	}
	return d.Dialector.DataTypeOf(field)
}

func (d SQLiteEnumDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return sqlite.Migrator{Migrator: migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   d,
		CreateIndexAfterCreateTable: true,
	}}}
}

// NewSQLiteEnumDialector creates a SQLite dialector that rewrites enum types for SQLite compatibility.
func NewSQLiteEnumDialector(dsn string) gorm.Dialector {
	return SQLiteEnumDialector{sqlite.Open(dsn)}
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
