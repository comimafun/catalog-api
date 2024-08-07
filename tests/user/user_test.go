package user_test

import (
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/user"
	test_helper "catalog-be/tests/test_helper"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var sleepTime = time.Millisecond * 5

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

func TestUser(t *testing.T) {
	ctx := context.Background()
	connURL, _ := test_helper.GetConnURL(t, ctx)
	db := test_helper.SetupDb(t, connURL)

	userRepo := user.NewUserRepo(db)
	service := user.NewUserService(userRepo)

	t.Run("Create new user", func(t *testing.T) {
		newUser := entity.User{
			Name:  "john doe",
			Email: "test@test.com",
		}
		createdUser, err := service.CreateOne(newUser)

		assert.Nil(t, err, "Error creating user")
		assert.NotNil(t, createdUser, "Created user should not be nil")

		// Validate that the created user has the expected values
		assert.Equal(t, newUser.Name, createdUser.Name)
		assert.Equal(t, newUser.Email, createdUser.Email)

		// Allow some time for DB operations to complete
		time.Sleep(sleepTime)
	})

	t.Run("Get user by email", func(t *testing.T) {
		email := "test@test.com"
		userByEmail, err := service.FindOneByEmail(email)

		assert.Nil(t, err, "Error finding user by email")
		assert.NotNil(t, userByEmail, "User should not be nil")

		assert.Equal(t, email, userByEmail.Email, "Email should match")
	})

	t.Run("Get user by ID", func(t *testing.T) {
		// Create a user first
		newUser := entity.User{
			Name:  "jane doe",
			Email: "jane@test.com",
		}
		createdUser, err := service.CreateOne(newUser)
		if err != nil {
			t.Fatalf("Failed to create user for ID lookup: %v", err)
		}

		userByID, err := service.FindOneByID(createdUser.ID)

		assert.Nil(t, err, "Error finding user by ID")
		assert.NotNil(t, userByID, "User should not be nil")

		assert.Equal(t, createdUser.ID, userByID.ID, "ID should match")
		assert.Equal(t, newUser.Email, userByID.Email, "Email should match")
	})
}
