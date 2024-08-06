package router

import (
	"catalog-be/internal/middlewares"
	"catalog-be/internal/modules/auth"
	"catalog-be/internal/modules/circle"
	"catalog-be/internal/modules/circle/circle_report"
	"catalog-be/internal/modules/circle/referral"
	"catalog-be/internal/modules/event"
	"catalog-be/internal/modules/fandom"
	"catalog-be/internal/modules/product"
	"catalog-be/internal/modules/upload"
	"catalog-be/internal/modules/work_type"

	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	auth           *auth.AuthHandler
	authMiddleware *middlewares.AuthMiddleware
	fandom         *fandom.FandomHandler
	workType       *work_type.WorkTypeHandler
	circle         *circle.CircleHandler
	event          *event.EventHandler
	upload         *upload.UploadHandler
	product        *product.ProductHandler
	referral       *referral.ReferralHandler
	circleReport   *circle_report.CircleReportHandler
}

func (h *HTTP) RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":     "UP",
			"request_id": c.Locals("requestid"),
		})
	})

	v1 := app.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.Get("/google", h.auth.GetAuthURL)
	auth.Get("/google/callback", h.auth.GetGoogleCallback)
	auth.Post("/google/callback", h.auth.PostGoogleCallback)
	auth.Get("/refresh", h.auth.GetGenerateNewTokenAndRefreshToken)
	auth.Get("/self", h.authMiddleware.Init, h.auth.GetSelf)
	auth.Post("/logout", h.authMiddleware.IfAuthed, h.auth.PostLogout)

	fandom := v1.Group("/fandom")
	fandom.Post("/", h.fandom.PostCreateOneFandom)
	fandom.Put("/:id", h.fandom.PutUpdateOneFandom)
	fandom.Delete("/:id", h.fandom.DeleteOneFandomByID)
	fandom.Get("/", h.fandom.GetPaginatedFandoms)

	workType := v1.Group("/worktype")
	// For admin only account
	workType.Post("/", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.workType.PostCreateOneWorkType)
	workType.Put("/:id", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.workType.PutUpdateOneWorkType)
	workType.Delete("/:id", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.workType.DeleteOneWorkTypeByID)

	workType.Get("/all", h.workType.GetAllWorkTypes)

	circle := v1.Group("/circle")
	circle.Post("/onboard", h.authMiddleware.Init, h.circle.PostOnboardNewCircle)
	circle.Patch("/:circleid", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.PatchUpdateOneCircleByCircleID)
	circle.Post("/:circleid/publish", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.PostPublishOrUnpublishCircle)

	circle.Get("/", h.authMiddleware.IfAuthed, h.circle.GetPaginatedCircles)
	circle.Get("/bookmarked", h.authMiddleware.Init, h.circle.GetPaginatedBookmarkedCircles)
	circle.Get("/:slug", h.authMiddleware.IfAuthed, h.circle.GetOneCricleByCircleSlug)

	circle.Get("/:circleid/referral", h.circle.GetCircleReferralByCirclceID)

	circle.Post("/:id/bookmark", h.authMiddleware.Init, h.circle.PostBookmarkCircleByCircleID)
	circle.Delete("/:id/bookmark", h.authMiddleware.Init, h.circle.DeleteBookmarkCircleByCircleID)

	circle.Get("/:id/product", h.product.GetAllProductByCircleID)
	circle.Post("/:id/product", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.product.CreateOneProductByCircleID)
	circle.Put("/:id/product/:productid", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.product.UpdateOneProductByCircleID)
	circle.Delete("/:id/product/:productid", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.product.DeleteOneProductByCircleIDAndProductID)

	circle.Put("/:circleid/event", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.PutUpdateAttendingEventByCircleID)
	circle.Delete("/:circleid/event", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.DeleteAttendingEventByCircleID)

	circle.Post("/:id/report", h.authMiddleware.Init, h.circleReport.CreateReportCircle)

	event := v1.Group("/event")
	event.Post("/", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.event.CreateOneEvent)
	event.Get("/", h.event.GetPaginatedEvents)

	upload := v1.Group("/upload")
	upload.Post("/image", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.upload.PostUploadImage)

	referral := v1.Group("/referral")
	referral.Post("/", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.referral.CreateOneReferral)
}

func NewHTTP(
	auth *auth.AuthHandler,
	authMiddleware *middlewares.AuthMiddleware,
	fandom *fandom.FandomHandler,
	workType *work_type.WorkTypeHandler,
	circle *circle.CircleHandler,
	event *event.EventHandler,
	upload *upload.UploadHandler,
	product *product.ProductHandler,
	referral *referral.ReferralHandler,
	circleReport *circle_report.CircleReportHandler,
) *HTTP {
	return &HTTP{
		auth,
		authMiddleware,
		fandom,
		workType,
		circle,
		event,
		upload,
		product,
		referral,
		circleReport,
	}
}
