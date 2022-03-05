package handler

import (
	"order-system/dao"
	"order-system/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Handler struct {
	userDAO    *dao.UserDAO
	productDAO *dao.ProductDAO
	orderDAO   *dao.OrderDAO
	validator  *utils.Validator
}

func NewHandler(us *dao.UserDAO, pd *dao.ProductDAO, od *dao.OrderDAO) *Handler {
	v := utils.NewValidator()
	return &Handler{
		userDAO:    us,
		productDAO: pd,
		orderDAO:   od,
		validator:  v,
	}
}

func (h *Handler) Register(r *fiber.App) {
	v1 := r.Group("/api")
	jwtMiddleware := jwtware.New(
		jwtware.Config{
			SigningKey: utils.JWTSecret,
			AuthScheme: "Token",
		})
	guestUsers := v1.Group("/users")
	guestUsers.Post("", h.SignUp)
	guestUsers.Post("/login", h.Login)
	user := v1.Group("/user", jwtMiddleware)
	user.Get("", h.CurrentUser)
	user.Put("", h.UpdateUser)

	productJWTMiddleware := jwtware.New(
		jwtware.Config{
			SigningKey: utils.JWTSecret,
			AuthScheme: "Token",
			Filter: func(c *fiber.Ctx) bool {
				if c.Method() == "GET" && strings.Contains(c.Path(), "/api/products") {
					return true
				}
				return false
			},
		})
	products := v1.Group("/products", productJWTMiddleware)
	products.Get("/:id", h.GetProduct)
	products.Get("", h.Products)
	products.Post("", h.CreateProduct)
	products.Put("/:id", h.UpdateProduct)

	carts := v1.Group("/carts", jwtMiddleware)
	carts.Get("", h.Carts)
	carts.Post("", h.CreateCartItem)
	carts.Put("/:id", h.UpdateCartItem)
	carts.Delete("/:id", h.DeleteCartItem)

	orders := v1.Group("/orders", jwtMiddleware)
	orders.Post("", h.PlaceOrder)
	orders.Get("/:id", h.GetOrder)
	orders.Get("", h.Orders)
}
