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
	if len(mdb.sqls) == 0 {
		fmt.Println("Nothing to generate.")
		return nil
	}
	filename := fmt.Sprintf("./migrations/%d_%s.up.sql", time.Now().UTC().Unix(), migrationName)
	os.Mkdir("./migrations/", os.ModeDir)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(mdb.GetMigrationString()))
	if err != nil {
		return err
	}
	return f.Close()
}
