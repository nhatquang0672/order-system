package handler

import (
	"order-system/model"
	"order-system/utils"

	"github.com/gofiber/fiber/v2"
)

type userUpdateRequest struct {
	User struct {
		Email    string `json:"email" validate:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"user"`
}

func newUserUpdateRequest() *userUpdateRequest {
	return new(userUpdateRequest)
}

func (r *userUpdateRequest) populate(u *model.User) {
	r.User.Email = u.Email
	r.User.Password = u.Password
	r.User.Username = u.Username
}

func (r *userUpdateRequest) bind(c *fiber.Ctx, u *model.User, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}
	if err := v.Validate(r); err != nil {
		return err
	}
	u.Email = r.User.Email
	//fmt.Printf("request user %v, from db user %v", r.User, *u)
	if r.User.Password != u.Password {
		h, err := u.HashPassword(r.User.Password)
		if err != nil {
			return err
		}
		u.Password = h
	}
	u.Username = r.User.Username
	return nil
}

type userRegisterRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	} `json:"user"`
}

func (r *userRegisterRequest) bind(c *fiber.Ctx, u *model.User, v *utils.Validator) error {
	//validate

	if err := c.BodyParser(r); err != nil {
		return err
	}
	//fmt.Printf("%v", *r)

	if err := v.Validate(r); err != nil {
		return err
	}
	u.Role = r.User.Role
	u.Email = r.User.Email
	h, err := u.HashPassword(r.User.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type userLoginRequest struct {
	User struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user"`
}

func (r *userLoginRequest) bind(c *fiber.Ctx, v *utils.Validator) error {

	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Validate(r); err != nil {
		return err
	}
	return nil
}

type productRequest struct {
	Product struct {
		Description string `json:"description"`
		Quantity    uint   `json:"quantity" validate:"numeric,gt=0"`
		Price       uint   `json:"price" validate:"numeric,gt=0"`
	} `json:"product"`
}

func (r *productRequest) populate(p *model.Product) {
	r.Product.Description = p.Description
	r.Product.Price = p.Price
	r.Product.Quantity = p.Quantity
}

func (r *productRequest) bind(c *fiber.Ctx, p *model.Product, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Validate(r); err != nil {
		return err
	}
	p.Description = r.Product.Description
	p.Price = r.Product.Price
	p.Quantity = r.Product.Quantity

	return nil
}

type cartItemRequest struct {
	CartItem struct {
		Quantity  uint `json:"quantity" validate:"numeric,gt=0"`
		ProductID uint `json:"productid"`
	} `json:"cartItem"`
}

func (r *cartItemRequest) populate(ct *model.CartItem) {
	r.CartItem.Quantity = ct.Quantity
	r.CartItem.ProductID = ct.ProductID
}

func (r *cartItemRequest) bind(c *fiber.Ctx, ct *model.CartItem, uid uint, v *utils.Validator) error {
	if err := c.BodyParser(r); err != nil {
		return err
	}

	if err := v.Validate(r); err != nil {
		return err
	}

	ct.Quantity = r.CartItem.Quantity
	ct.ProductID = r.CartItem.ProductID
	ct.UserID = uid
	return nil
}
