package handler

import (
	"order-system/model"
	"order-system/utils"
	"time"
)

type userResponse struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Token    string `json:"token"`
	} `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
	r := new(userResponse)
	r.User.Username = u.Username
	r.User.Email = u.Email
	r.User.Role = u.Role
	r.User.Token = utils.GenerateJWT(u.ID)
	return r
}

type productResponse struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Quantity    uint      `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
type singleProductResponse struct {
	Product *productResponse `json:"product"`
}
type listProductResponse struct {
	Products []productResponse `json:"products"`
	Count    int64             `json:"count"`
}

func newProductResponse(pd *model.Product) *singleProductResponse {
	product := new(productResponse)
	product.ID = pd.ID
	product.Description = pd.Description
	product.Price = pd.Price
	product.Quantity = pd.Quantity
	product.CreatedAt = pd.CreatedAt
	product.UpdatedAt = pd.UpdatedAt
	return &singleProductResponse{product}
}

func newListProductResponse(products []model.Product, count int64) *listProductResponse {
	r := new(listProductResponse)
	r.Count = count
	r.Products = make([]productResponse, 0)
	cr := productResponse{}
	for _, i := range products {
		cr.ID = i.ID
		cr.Description = i.Description
		cr.Price = i.Price
		cr.Quantity = i.Quantity
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		r.Products = append(r.Products, cr)
	}
	return r
}

type cartItemResponse struct {
	ID         uint      `json:"id"`
	Quantity   uint      `json:"quantity"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	TotalPrice uint      `json:"totalPrice"`
	Product    struct {
		ID          uint   `json:"id"`
		Description string `json:"description"`
		Price       uint   `json:"price"`
		Quantity    uint   `json:"quantity"`
	} `json:"product"`
}
type singleCartItemResponse struct {
	CartItem *cartItemResponse `json:"cartItem"`
}
type listCartItemsResponse struct {
	CartItems  []cartItemResponse `json:"cartItems"`
	TotalPrice uint               `json:"totalPrice"`
}

func newCartItemResponse(ct *model.CartItem) *singleCartItemResponse {
	cartItem := new(cartItemResponse)
	cartItem.ID = ct.ID
	cartItem.Quantity = ct.Quantity
	cartItem.CreatedAt = ct.CreatedAt
	cartItem.UpdatedAt = ct.UpdatedAt
	cartItem.Product.ID = ct.Product.ID
	cartItem.Product.Price = ct.Product.Price
	cartItem.Product.Description = ct.Product.Description
	cartItem.Product.Quantity = ct.Product.Quantity
	cartItem.TotalPrice = cartItem.Product.Price * cartItem.Quantity
	return &singleCartItemResponse{cartItem}
}

func newListCartItemResponse(cts []model.CartItem) *listCartItemsResponse {
	r := new(listCartItemsResponse)
	var totalPrice uint = 0
	r.CartItems = make([]cartItemResponse, 0)
	cr := cartItemResponse{}
	for _, i := range cts {
		cr.ID = i.ID
		cr.Quantity = i.Quantity
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		cr.Product.ID = i.Product.ID
		cr.Product.Price = i.Product.Price
		cr.Product.Description = i.Product.Description
		cr.Product.Quantity = i.Product.Quantity
		cr.TotalPrice = cr.Product.Price * cr.Quantity
		totalPrice += cr.TotalPrice
		r.CartItems = append(r.CartItems, cr)
	}
	r.TotalPrice = totalPrice
	return r
}

type orderItemResponse struct {
	ID        uint      `json:"id"`
	Quantity  uint      `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Price     uint      `json:"Price"`
	Product   struct {
		ID          uint   `json:"id"`
		Description string `json:"description"`
		Price       uint   `json:"price"`
		Quantity    uint   `json:"quantity"`
	} `json:"product"`
}
type orderRespone struct {
	Items      []orderItemResponse `json:"items"`
	TotalPrice uint                `json:"totalPrice"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
	ID         uint                `json:"id"`
}

type singleOrderResponse struct {
	Order *orderRespone `json:"order"`
}
type listOrderResponse struct {
	Orders []orderRespone `json:"orders"`
}

func newSingleOrderResponse(order *model.Order) *singleOrderResponse {
	r := new(orderRespone)
	r.ID = order.ID
	r.TotalPrice = order.TotalPrice
	r.Status = order.Status
	r.CreatedAt = order.CreatedAt
	r.UpdatedAt = order.UpdatedAt
	r.Items = make([]orderItemResponse, 0)
	cr := orderItemResponse{}
	for _, i := range order.Items {
		cr.ID = i.ID
		cr.Quantity = i.Quantity
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		cr.Product.ID = i.Product.ID
		cr.Product.Price = i.Product.Price
		cr.Product.Description = i.Product.Description
		cr.Product.Quantity = i.Product.Quantity
		cr.Price = i.Price
		r.Items = append(r.Items, cr)
	}
	return &singleOrderResponse{r}
}

func newListOrderResponse(orders []model.Order) *listOrderResponse {
	r := new(listOrderResponse)
	r.Orders = make([]orderRespone, 0)
	od := orderRespone{}
	for _, i := range orders {
		od.ID = i.ID
		od.TotalPrice = i.TotalPrice
		od.Status = i.Status
		od.CreatedAt = i.CreatedAt
		od.UpdatedAt = i.UpdatedAt
		od.Items = make([]orderItemResponse, 0)
		cr := orderItemResponse{}
		for _, j := range i.Items {
			cr.ID = i.ID
			cr.Quantity = j.Quantity
			cr.CreatedAt = j.CreatedAt
			cr.UpdatedAt = j.UpdatedAt
			cr.Product.ID = j.Product.ID
			cr.Product.Price = j.Product.Price
			cr.Product.Description = j.Product.Description
			cr.Product.Quantity = j.Product.Quantity
			cr.Price = j.Price
			od.Items = append(od.Items, cr)
		}
		r.Orders = append(r.Orders, od)
	}
	return r
}
