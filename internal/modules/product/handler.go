package product

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	auth_dto "catalog-be/internal/modules/auth/dto"
	product_dto "catalog-be/internal/modules/product/dto"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService *ProductService
	validator      *validator.Validate
}

func (p *ProductHandler) GetAllProductByCircleID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(400, err, nil)))
	}

	products, productErr := p.productService.GetAllProductsByCircleID(id)
	if productErr != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(domain.NewErrorFiber(c, productErr))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": products,
	})
}

func (p *ProductHandler) CreateOneProductByCircleID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)
	if *user.CircleID != id {
		return c.
			Status(fiber.StatusForbidden).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusForbidden, errors.New("FORBIDDEN"), nil)))
	}

	var body product_dto.CreateUpdateProductBody
	if err := c.BodyParser(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := p.validator.Struct(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	count, countErr := p.productService.CountProductsByCircleID(id)
	if countErr != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(domain.NewErrorFiber(c, countErr))
	}

	if count >= 5 {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("MAX_PRODUCT_EXCEEDED"), nil)))
	}

	product, productErr := p.productService.CreateOneProductByCircleID(id, entity.Product{
		Name:     body.Name,
		ImageURL: body.ImageURL,
	})
	if productErr != nil {
		return c.
			Status(productErr.Code).
			JSON(domain.NewErrorFiber(c, productErr))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": product,
	})
}

func (p *ProductHandler) UpdateOneProductByCircleID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	productID, err := c.ParamsInt("productid")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)
	if *user.CircleID != id {
		return c.
			Status(fiber.StatusForbidden).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusForbidden, errors.New("FORBIDDEN"), nil)))
	}

	var body product_dto.CreateUpdateProductBody
	if err := c.BodyParser(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := p.validator.Struct(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	product, productErr := p.productService.UpdateOneProductByCircleAndProductID(id, entity.Product{
		ID:       productID,
		Name:     body.Name,
		ImageURL: body.ImageURL,
	})
	if productErr != nil {
		return c.
			Status(productErr.Code).
			JSON(domain.NewErrorFiber(c, productErr))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": product,
	})
}

func (p *ProductHandler) DeleteOneProductByCircleIDAndProductID(c *fiber.Ctx) error {
	circleID, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	productID, err := c.ParamsInt("productid")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)
	if *user.CircleID != circleID {
		return c.
			Status(fiber.StatusForbidden).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusForbidden, errors.New("FORBIDDEN"), nil)))
	}

	deleteErr := p.productService.DeleteOneProductByID(circleID, productID)

	if deleteErr != nil {
		return c.
			Status(deleteErr.Code).
			JSON(domain.NewErrorFiber(c, deleteErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": "PRODUCT_DELETED",
	})

}

func NewProductHandler(productService *ProductService, validator *validator.Validate) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		validator:      validator,
	}
}
