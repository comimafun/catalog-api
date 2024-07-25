package router

import (
	"catalog-be/internal/middlewares"
	"catalog-be/internal/modules/auth"
	"catalog-be/internal/modules/circle"
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
	auth.Get("/refresh", h.auth.RefreshToken)
	auth.Get("/self", h.authMiddleware.Init, h.auth.GetSelf)
	auth.Post("/logout", h.authMiddleware.IfAuthed, h.auth.Logout)

	fandom := v1.Group("/fandom")
	fandom.Post("/", h.fandom.CreateOne)
	fandom.Put("/:id", h.fandom.UpdateOne)
	fandom.Delete("/:id", h.fandom.DeleteByID)
	fandom.Get("/", h.fandom.GetFandomPagination)

	workType := v1.Group("/worktype")
	// For admin only account
	workType.Post("/", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.workType.CreateOne)
	workType.Put("/:id", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.workType.UpdateOne)
	workType.Delete("/:id", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.workType.DeleteByID)

	workType.Get("/all", h.workType.GetAll)

	circle := v1.Group("/circle")
	circle.Post("/onboard", h.authMiddleware.Init, h.circle.OnboardNewCircle)
	circle.Patch("/:circleid", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.UpdateCircle)
	circle.Post("/:circleid/publish", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.PublishUnpublishCircle)

	circle.Get("/", h.authMiddleware.IfAuthed, h.circle.GetPaginatedCircle)
	circle.Get("/bookmarked", h.authMiddleware.Init, h.circle.GetPaginatedBookmarkedCircle)
	circle.Get("/:slug", h.authMiddleware.IfAuthed, h.circle.FindCircleBySlug)

	circle.Post("/:id/bookmark", h.authMiddleware.Init, h.circle.SaveCircle)
	circle.Delete("/:id/bookmark", h.authMiddleware.Init, h.circle.UnsaveCircle)

	circle.Get("/:id/product", h.product.GetAllProductByCircleID)
	circle.Post("/:id/product", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.product.CreateOneProduct)
	circle.Put("/:id/product/:productid", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.product.UpdateOneProduct)
	circle.Delete("/:id/product/:productid", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.product.DeleteOneProduct)

	circle.Put("/:id/event", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.UpdateCircleEventAttending)
	circle.Delete("/:id/event", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.circle.DeleteCircleEventAttending)

	event := v1.Group("/event")
	event.Post("/", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.event.CreateOne)
	event.Get("/", h.event.GetPaginatedEvents)

	upload := v1.Group("/upload")
	upload.Post("/image", h.authMiddleware.Init, h.authMiddleware.CircleOnly, h.upload.UploadImage)

	referral := v1.Group("/referral")
	referral.Post("/", h.authMiddleware.Init, h.authMiddleware.AdminOnly, h.referral.CreateReferral)
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
	}
}
