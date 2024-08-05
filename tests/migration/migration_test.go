package migration_test

import (
	"catalog-be/internal/database"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

func getConnUrl(t *testing.T, ctx context.Context) string {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("foobar"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
	}

	connUrl, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Could not get connection string: %s", err)
	}

	return connUrl
}

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
	connUrl := getConnUrl(t, ctx)
	db := database.New(connUrl, false)

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
