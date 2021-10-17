package order

type OrderRequest struct {
	Id string  `json:",omitempty"`
	CustomerId int `json:"customerId" validate:"required"`
	ProductId int `json:"productId" validate:"required"`
	OrderDesc string `json:"orderDesc" validate:"required"`
	ReceiptHandle *string `json:"-"`
	NumberOfOrders int
}

type OrderResponse struct {
	Code   string `json:"code"`
	Desc   string `json:"err,omitempty"`
}