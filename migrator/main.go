package migrator

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func autoMigrate(db *gorm.DB, models ...interface{}) (strs []string, err error) {
	db.Callback().Raw().Remove("gorm:raw")
	db.Callback().Raw().Register("migrations", func(db *gorm.DB) {
		stmt := db.Statement
		sql := stmt.SQL.String()
		if strings.HasPrefix(sql, "CREATE") || strings.HasPrefix(sql, "ALTER") || strings.HasPrefix(sql, "DROP") {
			strs = append(strs, sql)
		} else if db.Error == nil && !db.DryRun {
			result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				db.AddError(err)
			} else {
				db.RowsAffected, _ = result.RowsAffected()
			}
		}
	})
	err = db.AutoMigrate(models...)
	db.Callback().Raw().Remove("migrations")
	db.Callback().Raw().Register("gorm:raw", callbacks.RawExec)
	return
}

// GenerateMigrations writes migration files.
func GenerateMigrations(migrationName string, db *gorm.DB, values ...interface{}) error {
	sqls, err := autoMigrate(db, values...)
	if err != nil {
		return err
	}
	if len(sqls) == 0 {
		fmt.Println("Nothing to generate.")
		return nil
	}
	filename := fmt.Sprintf("./migrations/%d_%s.up.sql", time.Now().UTC().Unix(), migrationName)
	os.Mkdir("./migrations/", os.ModeDir)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(strings.Join(sqls, ";\n")))
	if err != nil {
		return err
	}
	return f.Close()
}
