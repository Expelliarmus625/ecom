package products

import (
	"context"
	"log"
	"strconv"

	"github.com/expelliarmus625/ecom/internal/adapters/postgresql/sqlc"
)

type Service interface{
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductByID(ctx context.Context, productId string) (repo.Product, error)
}

type svc struct{
	repo repo.Querier
}

func NewService(repo repo.Querier) Service{
	return &svc{repo: repo}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductByID(ctx context.Context, productId string) (repo.Product, error) {
	id, err := strconv.Atoi(productId)
	if err != nil {
		return repo.Product{}, err
	}
	
	log.Default().Printf("%d", id)
	return s.repo.FindProductByID(ctx, int64(id))
}
