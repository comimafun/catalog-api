//go:build wireinject
// +build wireinject

package internal

import (
	internal_config "catalog-be/internal/config"
	"catalog-be/internal/middlewares"
	"catalog-be/internal/modules/auth"
	"catalog-be/internal/modules/circle"
	"catalog-be/internal/modules/circle/bookmark"
	"catalog-be/internal/modules/circle/circle_fandom"
	"catalog-be/internal/modules/circle/circle_work_type"
	"catalog-be/internal/modules/event"
	"catalog-be/internal/modules/fandom"
	"catalog-be/internal/modules/product"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/modules/work_type"
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

		fandom.NewFandomHandler,
		fandom.NewFandomRepo,
		fandom.NewFandomService,

		work_type.NewWorkTypeHandler,
		work_type.NewWorkTypeRepo,
		work_type.NewWorkTypeService,

		bookmark.NewCircleBookmarkRepo,
		bookmark.NewCircleBookmarkService,

		circle_work_type.NewCircleWorkTypeRepo,
		circle_work_type.NewCircleWorkTypeService,

		circle_fandom.NewCircleFandomRepo,
		circle_fandom.NewCircleFandomService,

		// circle_upvote.NewCircleUpvoteRepo,
		// circle_upvote.NewCircleUpvoteService,

		product.NewProductRepo,
		product.NewProductService,

		circle.NewCircleHandler,
		circle.NewCircleRepo,
		circle.NewCircleService,

		event.NewEventHandler,
		event.NewEventRepo,
		event.NewEventService,

		middlewares.NewAuthMiddleware,

		router.NewHTTP,
	)
	return nil
}
