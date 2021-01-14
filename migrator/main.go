package migrator

import (
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

// GenerateMigrations writes migration files.
func GenerateMigrations(migrationName string, db *gorm.DB, values ...interface{}) error {
	mdb := NewMigratorDB(*db)
	err := mdb.AutoMigrate(values)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("./migrations/%s_%s.up.sql", time.Now().UTC().Format(time.Stamp), migrationName)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(mdb.GetMigrationString()))
	if err != nil {
		return err
	}
	return f.Close()
}
