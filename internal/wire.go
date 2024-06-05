//go:build wireinject
// +build wireinject

package internal

import (
	internal_config "catalog-be/internal/config"
	"catalog-be/internal/middlewares"
	"catalog-be/internal/modules/auth"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/router"
	"catalog-be/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeServer(db *gorm.DB, validator *validator.Validate) *router.HTTP {
	wire.Build(
		internal_config.NewConfig,

		utils.NewUtils,

		refreshtoken.NewRefreshTokenRepo,
		refreshtoken.NewRefreshTokenService,

		user.NewUserRepo,
		user.NewUserService,

		auth.NewAuthHandler,
		auth.NewAuthService,

		middlewares.NewAuthMiddleware,

		router.NewHTTP,
	)
	return nil
}
