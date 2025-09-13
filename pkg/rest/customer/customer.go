package customer

import (
	"context"
	"fmt"

	client "github.com/rasteiro11/MCABankGateway/pkg/transport/http"
	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
)

type customerClient struct {
	BaseURL string
}

var _ CustomerClient = (*customerClient)(nil)

func New(baseURL string) CustomerClient {
	return &customerClient{BaseURL: baseURL}
}

func ToCustomerPointers(customers []domain.Customer) []*domain.Customer {
	result := make([]*domain.Customer, len(customers))
	for i := range customers {
		result[i] = &customers[i]
	}
	return result
}

func (c *customerClient) GetAll(ctx context.Context, queryParams map[string]string) ([]*domain.Customer, error) {
	path := fmt.Sprintf("%s/customers", c.BaseURL)
	opts := []client.Option{}
	if len(queryParams) > 0 {
		opts = append(opts, client.WithQueryParams(queryParams))
	}

	res, err := client.Get[[]domain.Customer](ctx, path, opts...)
	if err != nil {
		return nil, err
	}

	return ToCustomerPointers(*res), nil
}

func (c *customerClient) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {
	path := fmt.Sprintf("%s/customers/%d", c.BaseURL, id)
	return client.Get[domain.Customer](ctx, path)
}

func (c *customerClient) Create(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	path := fmt.Sprintf("%s/customers", c.BaseURL)
	return client.Post[domain.Customer](ctx, path, customer)
}

func (c *customerClient) Update(ctx context.Context, id uint, customer *domain.Customer) (*domain.Customer, error) {
	path := fmt.Sprintf("%s/customers/%d", c.BaseURL, id)
	return client.Put[domain.Customer](ctx, path, customer)
}

func (c *customerClient) Delete(ctx context.Context, id uint) error {
	path := fmt.Sprintf("%s/customers/%d", c.BaseURL, id)
	_, err := client.Delete[any](ctx, path)
	return err
}
