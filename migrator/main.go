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
	err := mdb.AutoMigrate(values...)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("./migrations/%d_%s.up.sql", time.Now().UTC().Unix(), migrationName)
	os.Mkdir("./migrations/", os.ModeDir)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(mdb.GetMigrationString()))
	if err != nil {
		return err
	}
	return f.Close()
}
