package orders

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

func TestNewServer(t *testing.T) {
	t.Run("new_server_creates_handlers", func(t *testing.T) {
		mockClient := &mockOrderServiceClient{}
		server := NewServer(mockClient)

		if server.createOrderHandler == nil {
			t.Error("createOrderHandler should not be nil")
		}

		if server.getOrderHandler == nil {
			t.Error("getOrderHandler should not be nil")
		}

		if server.listOrdersHandler == nil {
			t.Error("listOrdersHandler should not be nil")
		}

		if server.client != mockClient {
			t.Error("client should be set")
		}
	})
}

func TestServer_CheckOrderOwner(t *testing.T) {
	t.Run("check_order_owner_unimplemented", func(t *testing.T) {
		mockClient := &mockOrderServiceClient{}
		server := NewServer(mockClient)

		_, err := server.CheckOrderOwner(context.Background(), connect.NewRequest(&orderv1.CheckOrderOwnerRequest{
			OrderId: "order-123",
			UserId:  "user-123",
		}))

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		connectErr, ok := err.(*connect.Error)
		if !ok {
			t.Fatalf("expected connect.Error, got %T", err)
		}

		if connectErr.Code() != connect.CodeUnimplemented {
			t.Errorf("expected Unimplemented, got %v", connectErr.Code())
		}
	})
}

func TestServer_Delegation(t *testing.T) {
	t.Run("delegates_to_create_order_handler", func(t *testing.T) {
		expectedOrder := &orderv1.Order{
			Id:     "order-123",
			UserId: "user-123",
			Item:   "test-item",
			Amount: 100,
		}

		mockClient := &mockOrderServiceClient{
			createOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
				return connect.NewResponse(&orderv1.CreateOrderResponse{
					Order: expectedOrder,
				}), nil
			},
		}

		server := NewServer(mockClient)
		resp, err := server.CreateOrder(context.Background(), connect.NewRequest(&orderv1.CreateOrderRequest{
			UserId: "user-123",
			Item:   "test-item",
			Amount: 100,
		}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Msg.Order.Id != expectedOrder.Id {
			t.Errorf("expected order ID %s, got %s", expectedOrder.Id, resp.Msg.Order.Id)
		}
	})

	t.Run("delegates_to_get_order_handler", func(t *testing.T) {
		expectedOrder := &orderv1.Order{
			Id:     "order-123",
			UserId: "user-123",
			Item:   "test-item",
			Amount: 100,
		}

		mockClient := &mockOrderServiceClient{
			checkOrderOwnerFunc: func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
				return connect.NewResponse(&orderv1.CheckOrderOwnerResponse{}), nil
			},
			getOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
				return connect.NewResponse(&orderv1.GetOrderResponse{
					Order: expectedOrder,
				}), nil
			},
		}

		server := NewServer(mockClient)
		resp, err := server.GetOrder(context.Background(), connect.NewRequest(&orderv1.GetOrderRequest{
			Id:     "order-123",
			UserId: "user-123",
		}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Msg.Order.Id != expectedOrder.Id {
			t.Errorf("expected order ID %s, got %s", expectedOrder.Id, resp.Msg.Order.Id)
		}
	})

	t.Run("delegates_to_list_orders_handler", func(t *testing.T) {
		expectedOrders := []*orderv1.Order{
			{
				Id:     "order-1",
				UserId: "user-123",
				Item:   "item-1",
				Amount: 100,
			},
			{
				Id:     "order-2",
				UserId: "user-123",
				Item:   "item-2",
				Amount: 200,
			},
		}

		mockClient := &mockOrderServiceClient{
			listOrdersFunc: func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
				return connect.NewResponse(&orderv1.ListOrdersResponse{
					Orders: expectedOrders,
				}), nil
			},
		}

		server := NewServer(mockClient)
		resp, err := server.ListOrders(context.Background(), connect.NewRequest(&orderv1.ListOrdersRequest{}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(resp.Msg.Orders) != len(expectedOrders) {
			t.Errorf("expected %d orders, got %d", len(expectedOrders), len(resp.Msg.Orders))
		}
	})
}
