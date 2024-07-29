package user_test

import (
	"catalog-be/internal/database"
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/user"
	"context"
	"fmt"
	"log"
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

var connUrl = ""
var sleepTime = time.Millisecond * 5

func migrate(db *gorm.DB) {
	sqlDir := "../migrator/migrations"

	files, err := os.ReadDir(sqlDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "up.sql") {
			sqlContent, err := os.ReadFile(fmt.Sprintf("%s/%s", sqlDir, file.Name()))
			if err != nil {
				panic(err)
			}

			err = db.Exec(string(sqlContent)).Error
			if err != nil {
				panic(err)
			}

			fmt.Printf("Migrated: %s\n", file.Name())
		}
	}
}

func setupDb(dsn string) *gorm.DB {
	db := database.New(dsn)
	migrate(db)
	return db
}

func TestMain(m *testing.M) {
	ctx := context.Background()

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
		log.Fatalf("Could not start postgres container: %s", err)
	}

	connUrl, err = container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Could not get connection string: %s", err)
	}
	setupDb(connUrl)
	res := m.Run()
	os.Exit(res)
}

func TestUser(t *testing.T) {
	db := database.New(connUrl)
	userRepo := user.NewUserRepo(db)

	service := user.NewUserService(userRepo)

	t.Run("Create new user", func(t *testing.T) {
		as, err := service.CreateOne(entity.User{
			Name:  "john doe",
			Email: "test@test.com",
		})

		assert.Nil(t, err)
		assert.NotNil(t, as)

		time.Sleep(sleepTime)
	})

	t.Run("Get user by email", func(t *testing.T) {
		as, err := service.FindOneByEmail("test@test.com")

		assert.Nil(t, err)
		assert.NotNil(t, as)

		assert.Equal(t, "test@test.com", as.Email)
		assert.Equal(t, 1, as.ID)
	})
}
