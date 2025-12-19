package orders

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	repo "github.com/expelliarmus625/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock = errors.New("product not enough stock")
)

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	ListOrders(ctx context.Context) ([]repo.Order, error)
	ListOrderItems(ctx context.Context, id string) ([]repo.OrderItem, error)
}

type svc struct {
	repo *repo.Queries
	db *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db: 	db,
	}
}

func (s *svc) ListOrderItems(ctx context.Context, id string) ([]repo.OrderItem, error) {
	orderId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return []repo.OrderItem{}, err
	}
	orderItems, err := s.repo.ListOrderItems(ctx, orderId)

	if len(orderItems) == 0 {
		return []repo.OrderItem{}, fmt.Errorf("order id not found")
	}

	return orderItems, err
}

func (s *svc) ListOrders(ctx context.Context) ([]repo.Order, error) {
	return s.repo.ListOrders(ctx)
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error){
	//validate the payload
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("CustomerID cannot be null")
	}

	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("Atleast one order item required")
	}

	//create an order
	tx, err := s.db.Begin(ctx)
	if err != nil{
		return repo.Order{}, err
	}	

	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}

	//look for product if exists
	for _, item := range tempOrder.Items {
		product, err := s.repo.FindProductByID(ctx, item.ProductID)
		if err != nil{
			return repo.Order{}, ErrProductNotFound
		}

		if product.Quantity < int32(item.Quantity){
			return repo.Order{}, ErrProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID: order.ID,
			ProductID: item.ProductID,
			Quantity: int32(item.Quantity),
			PriceCents: product.PriceInCents,
		})

		if err != nil {
			return repo.Order{}, err
		}
	}
	 tx.Commit(ctx)

	return order, nil
}
