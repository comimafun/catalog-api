package product

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	product_dto "catalog-be/internal/modules/product/dto"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService ProductService
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

func (p *ProductHandler) CreateOneProduct(c *fiber.Ctx) error {
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

	var body product_dto.CreateProductBody
	if err := c.BodyParser(body); err != nil {
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

	product, productErr := p.productService.UpsertOneProductByCircleID(id, product_dto.UpdateProductBody{
		CreateProductBody: body,
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

func (p *ProductHandler) UpdateOneProduct(c *fiber.Ctx) error {
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

	var body product_dto.UpdateProductBody
	if err := c.BodyParser(body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	product, productErr := p.productService.UpsertOneProductByCircleID(id, body)
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

func (p *ProductHandler) UpsertCircleProducts(c *fiber.Ctx) error {
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

	var body []product_dto.UpdateProductBody
	if err := c.BodyParser(body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if len(body) > 5 {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("MAX_PRODUCT_EXCEEDED"), nil)))
	}

	products, productErr := p.productService.UpsertProductByCircleID(id, body)
	if productErr != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(domain.NewErrorFiber(c, productErr))
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": products,
	})
}

func NewProductHandler(productService ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}
