package orders

type orderItem struct{
	ProductID int64 `json:"product_id"`
	Quantity int `json:"quantity"`
}

type createOrderParams struct{
	CustomerID int64 `json:"customer_id"`
	Items []orderItem `json:"items"`
}
