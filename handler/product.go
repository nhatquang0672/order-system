package handler

import (
	"errors"
	"log"
	"net/http"
	"order-system/model"
	"order-system/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetProduct godoc
// @Summary Get an product
// @Description Get an product. Auth not required
// @ID get-product
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path integer true "ID of the product you want to get"
// @Success 200 {object} singleProductResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /products/{id} [get]
func (h *Handler) GetProduct(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(err))
	}
	pd, err := h.productDAO.GetById(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if pd == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	return c.Status(http.StatusOK).JSON(newProductResponse(pd))
}

// Products godoc
// @Summary Get recent products globally
// @Description Get most recent products globally. Use query parameters to filter results. Auth is optional
// @ID get-products
// @Tags product
// @Accept json
// @Produce json
// @Param limit query integer false "Limit number of product returned (default is 20)"
// @Param offset query integer false "Offset/skip number of product (default is 0)"
// @Success 200 {object} listProductResponse
// @Failure 500 {object} utils.Error
// @Router /products [get]
func (h *Handler) Products(c *fiber.Ctx) error {
	var (
		products []model.Product
		count    int64
	)
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}
	products, count, err = h.productDAO.List(offset, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newListProductResponse(products, count))
}

// CreateProduct godoc
// @Summary Create new product
// @Description Create new product
// @ID create
// @Tags product
// @Accept json
// @Produce json
// @Param product body productRequest true "Product details to create. At least **one** field is required."
// @Success 200 {object} singleProductResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /products [post]
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	if u.Role != "vendor" {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}
	var pd model.Product

	req := &productRequest{}
	if err := req.bind(c, &pd, h.validator); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	if err := h.productDAO.Create(&pd); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newProductResponse(&pd))
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product information from product id.
// @ID update
// @Tags product
// @Accept json
// @Produce json
// @Param product body productRequest true "Product details to update. At least **one** field is required."
// @Param id path integer true "ID of the product you want to update"
// @Success 200 {object} singleProductResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	if u.Role != "vendor" {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}

	id64, err := strconv.ParseUint(c.Params("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(err))
	}
	pd, err := h.productDAO.GetById(id)

	if err != nil {
		log.Println("[FATAL]")
		log.Fatal(err)
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}

	if pd == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	req := &productRequest{}
	req.populate(pd)
	if err := req.bind(c, pd, h.validator); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	if err := h.productDAO.Update(pd); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newProductResponse(pd))
}
