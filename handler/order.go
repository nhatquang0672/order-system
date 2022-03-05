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

// Carts godoc
// @Summary Get cart items for user
// @Description Get all items in user's cart.
// @ID get-cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} listCartItemsResponse
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /carts [get]
func (h *Handler) Carts(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	var cartItems []model.CartItem

	cartItems, err = h.orderDAO.ListCartItems(u.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newListCartItemResponse(cartItems))
}

// CreateCartItem godoc
// @Summary Create cart item
// @Description Create cart item
// @ID create-cart-item
// @Tags cart
// @Accept json
// @Produce json
// @Param cartItem body cartItemRequest true "Cart item"
// @Success 200 {object} singleCartItemResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /carts [post]
func (h *Handler) CreateCartItem(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}
	var ct model.CartItem

	req := &cartItemRequest{}
	if err := req.bind(c, &ct, u.ID, h.validator); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	pd, err := h.productDAO.GetById(req.CartItem.ProductID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if pd == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}

	if pd.Quantity < req.CartItem.Quantity {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("exceed quantity")))
	}

	if err := h.orderDAO.CreateCartItem(&ct); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newCartItemResponse(&ct))
}

// UpdateCartItem godoc
// @Summary Update cart item
// @Description Update cart item
// @ID update-cart-item
// @Tags cart
// @Accept json
// @Produce json
// @Param cartItem body cartItemRequest true "Cart item"
// @Param id path integer true "ID of the cart item you want to update"
// @Success 200 {object} singleCartItemResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /carts/{id} [put]
func (h *Handler) UpdateCartItem(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}
	id64, err := strconv.ParseUint(c.Params("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(err))
	}

	cartItem, err := h.orderDAO.CartItemFromId(id)
	if err != nil {
		log.Println("[FATAL]")
		log.Fatal(err)
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if cartItem == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}

	req := &cartItemRequest{}
	req.populate(cartItem)

	if err := req.bind(c, cartItem, u.ID, h.validator); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}

	if req.CartItem.Quantity > cartItem.Product.Quantity {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("exceed quantity")))
	}

	if err := h.orderDAO.UpdateCartItem(cartItem); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newCartItemResponse(cartItem))
}

// DeleteCartItem godoc
// @Summary Delete cart item
// @Description Delete cart item
// @ID delete-cart-item
// @Tags cart
// @Accept json
// @Produce json
// @Param id path integer true "ID of the cart item you want to get"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /carts/{id} [delete]
func (h *Handler) DeleteCartItem(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}
	id64, err := strconv.ParseUint(c.Params("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(err))
	}
	cartItem, err := h.orderDAO.CartItemFromId(id)
	if err != nil {
		log.Println("[FATAL]")
		log.Fatal(err)
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if cartItem == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	if cartItem.UserID != u.ID {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))

	}
	if err := h.orderDAO.DeleteCartItem(cartItem); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(map[string]interface{}{"result": "ok"})
}

// PlaceOrder godoc
// @Summary Place order for user
// @Description Place order
// @ID place-order
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {object} singleOrderResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /orders [post]
func (h *Handler) PlaceOrder(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}
	var cartItems []model.CartItem
	cartItems, err = h.orderDAO.ListCartItems(u.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if len(cartItems) == 0 {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(errors.New("Your cart is empty!")))
	}

	order, err := h.orderDAO.PlaceOrder(u.ID, &cartItems)
	if err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(utils.NewError(err))
	}
	return c.Status(http.StatusOK).JSON(newSingleOrderResponse(order))
}

// GetOrder godoc
// @Summary Get order from order id
// @Description Get order
// @ID get-order
// @Tags order
// @Accept json
// @Produce json
// @Param id path integer true "ID of the order you want to get"
// @Success 200 {object} singleOrderResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}

	id64, err := strconv.ParseUint(c.Params("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.NewError(err))
	}

	order, err := h.orderDAO.GetOrderById(id)
	if order == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if order.UserID != u.ID {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}

	return c.Status(http.StatusOK).JSON(newSingleOrderResponse(order))
}

// Orders godoc
// @Summary Get list orders
// @Description Get list orders
// @ID get-list-orders
// @Tags order
// @Accept json
// @Produce json
// @Param limit query integer false "Limit number of order returned (default is 20)"
// @Param offset query integer false "Offset/skip number of order (default is 0)"
// @Success 200 {object} listOrderResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /orders [get]
func (h *Handler) Orders(c *fiber.Ctx) error {
	u, err := h.userDAO.GetByID(userIDFromToken(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if u == nil {
		return c.Status(http.StatusUnauthorized).JSON(utils.NewError(errors.New("unathorized action")))
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}

	orders, err := h.orderDAO.Orders(u.ID, offset, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(utils.NewError(err))
	}
	if orders == nil {
		return c.Status(http.StatusNotFound).JSON(utils.NotFound())
	}
	if len(orders) == 0 {
		return c.Status(http.StatusBadGateway).JSON(utils.NewError(errors.New("You don't have any order!")))
	}

	return c.Status(http.StatusOK).JSON(newListOrderResponse(orders))
}
