package mock

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

// OrderService is a mock implementation of OrderServiceHandler for testing.
type OrderService struct{}

// NewOrderService creates a new mock OrderService.
func NewOrderService() *OrderService {
	return &OrderService{}
}

// CreateOrder returns a mock order based on the request.
func (s *OrderService) CreateOrder(
	ctx context.Context,
	req *connect.Request[orderv1.CreateOrderRequest],
) (*connect.Response[orderv1.CreateOrderResponse], error) {
	orderID := fmt.Sprintf("order-%s-001", req.Msg.UserId)

	return connect.NewResponse(&orderv1.CreateOrderResponse{
		Order: &orderv1.Order{
			Id:     orderID,
			UserId: req.Msg.UserId,
			Item:   req.Msg.Item,
			Amount: req.Msg.Amount,
		},
	}), nil
}

// GetOrder returns a mock order by ID.
func (s *OrderService) GetOrder(
	ctx context.Context,
	req *connect.Request[orderv1.GetOrderRequest],
) (*connect.Response[orderv1.GetOrderResponse], error) {
	if req.Msg.Id == "not-found" {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("order not found"))
	}

	return connect.NewResponse(&orderv1.GetOrderResponse{
		Order: &orderv1.Order{
			Id:     req.Msg.Id,
			UserId: req.Msg.UserId,
			Item:   "mock-item",
			Amount: 100,
		},
	}), nil
}

// ListOrders returns a predefined list of orders.
func (s *OrderService) ListOrders(
	ctx context.Context,
	req *connect.Request[orderv1.ListOrdersRequest],
) (*connect.Response[orderv1.ListOrdersResponse], error) {
	return connect.NewResponse(&orderv1.ListOrdersResponse{
		Orders: []*orderv1.Order{
			{
				Id:     "order-001",
				UserId: "user-1",
				Item:   "item-1",
				Amount: 100,
			},
			{
				Id:     "order-002",
				UserId: "user-1",
				Item:   "item-2",
				Amount: 200,
			},
		},
	}), nil
}

// CheckOrderOwner validates order ownership.
func (s *OrderService) CheckOrderOwner(
	ctx context.Context,
	req *connect.Request[orderv1.CheckOrderOwnerRequest],
) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
	if req.Msg.UserId == "unauthorized" {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("user is not the owner"))
	}

	return connect.NewResponse(&orderv1.CheckOrderOwnerResponse{}), nil
}
