package migration_test

import (
	"catalog-be/internal/database"
	test_helper "catalog-be/tests/test_helper"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func migrateUp(t *testing.T, db *gorm.DB) error {
	fmt.Printf("[=] MIGRATING UP START [=]\n")
	sqlDir := "../../migrator/migrations"

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		t.Fatalf("Error reading directory: %s", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "up.sql") {
			sqlContent, err := os.ReadFile(fmt.Sprintf("%s/%s", sqlDir, file.Name()))
			if err != nil {
				t.Logf("PATH: %s/%s", sqlDir, file.Name())
				return err
			}

			err = db.Exec(string(sqlContent)).Error
			if err != nil {
				t.Logf("PATH: %s/%s", sqlDir, file.Name())
				return err
			}
		}
	}

	fmt.Printf("[=] MIGRATING UP END [=]\n")

	return nil
}

func migrateDown(t *testing.T, db *gorm.DB) error {
	fmt.Printf("[=] MIGRATING DOWN START [=]\n")
	sqlDir := "../../migrator/migrations"

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		return err
	}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		if strings.HasSuffix(file.Name(), "down.sql") {
			sqlContent, err := os.ReadFile(fmt.Sprintf("%s/%s", sqlDir, file.Name()))
			if err != nil {
				t.Logf("PATH: %s/%s", sqlDir, file.Name())
				return err
			}

			err = db.Exec(string(sqlContent)).Error
			if err != nil {
				t.Logf("PATH: %s/%s", sqlDir, file.Name())
				return err
			}
		}
	}

	fmt.Printf("[=] MIGRATING DOWN END [=]\n")

	return nil
}

func TestMigrationScript(t *testing.T) {
	ctx := context.Background()
	connURL, _ := test_helper.GetConnURL(t, ctx)
	db := database.New(connURL, false)

	err := migrateUp(t, db)
	if err != nil {
		t.Logf("Error migrating up: %s", err)
	}
	assert.Nil(t, err)

	err = migrateDown(t, db)
	if err != nil {
		t.Logf("Error migrating down: %s", err)
	}
	assert.Nil(t, err)
}
