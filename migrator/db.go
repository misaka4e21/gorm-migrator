package migrator

import (
	"strings"

	"gorm.io/gorm"
)

// DB is the gorm.DB subtype for gorm-migrator.
type DB struct {
	gorm.DB
	sqls []string
}

// Raw prevents CREATE, ALTER and DROP statements to be run, then records them.
func (db *DB) Raw(sql string, values ...interface{}) (tx *DB) {
	if strings.HasPrefix(sql, "CREATE") || strings.HasPrefix(sql, "ALTER") || strings.HasPrefix(sql, "DROP") {
		db.sqls = append(db.sqls, sql)
	} else {
		db.DB = *db.DB.Raw(sql, values...)
	}
	return db
}

// GetMigrationString returns the intercepted SQL statements.
func (db *DB) GetMigrationString() string {
	return strings.Join(db.sqls, ";\n")
}

// NewMigratorDB generates a new migrator.DB instance.
func NewMigratorDB(gormdb gorm.DB) *DB {
	return &DB{
		DB:   gormdb,
		sqls: []string{},
	}
}
