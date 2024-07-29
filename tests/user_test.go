package tests

import (
	"catalog-be/internal/database"
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	db := database.New(connUrl)
	userRepo := user.NewUserRepo(db)

	as, err := userRepo.CreateOne(entity.User{
		Name: "john doe",
	})

	assert.Nil(t, err)
	assert.NotNil(t, as)

	time.Sleep(sleepTime)
}
