package orders

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	orderv1 "github.com/demo/contracts/gen/go/order/v1"
	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
)

// mockOrderServiceClient is a mock implementation of orderv1connect.OrderServiceClient
// for unit testing handlers in isolation.
type mockOrderServiceClient struct {
	createOrderFunc     func(context.Context, *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error)
	getOrderFunc        func(context.Context, *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error)
	listOrdersFunc      func(context.Context, *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error)
	checkOrderOwnerFunc func(context.Context, *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error)
}

// Compile-time check that mockOrderServiceClient implements orderv1connect.OrderServiceClient.
var _ orderv1connect.OrderServiceClient = (*mockOrderServiceClient)(nil)

func (m *mockOrderServiceClient) CreateOrder(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, req)
	}
	return nil, errors.New("createOrderFunc not implemented")
}

func (m *mockOrderServiceClient) GetOrder(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
	if m.getOrderFunc != nil {
		return m.getOrderFunc(ctx, req)
	}
	return nil, errors.New("getOrderFunc not implemented")
}

func (m *mockOrderServiceClient) ListOrders(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
	if m.listOrdersFunc != nil {
		return m.listOrdersFunc(ctx, req)
	}
	return nil, errors.New("listOrdersFunc not implemented")
}

func (m *mockOrderServiceClient) CheckOrderOwner(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
	if m.checkOrderOwnerFunc != nil {
		return m.checkOrderOwnerFunc(ctx, req)
	}
	return nil, errors.New("checkOrderOwnerFunc not implemented")
}
