package orders

import (
	"context"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
)

type Server struct {
	client             orderv1connect.OrderServiceClient
	createOrderHandler *createOrderHandler
	getOrderHandler    *getOrderHandler
	listOrdersHandler  *listOrdersHandler
}

func NewServer(client orderv1connect.OrderServiceClient) *Server {
	return &Server{
		client:             client,
		createOrderHandler: newCreateOrderHandler(client),
		getOrderHandler:    newGetOrderHandler(client),
		listOrdersHandler:  newListOrdersHandler(client),
	}
}

func (s *Server) CreateOrder(
	ctx context.Context,
	req *connect.Request[orderv1.CreateOrderRequest],
) (*connect.Response[orderv1.CreateOrderResponse], error) {
	return s.createOrderHandler.Handle(ctx, req)
}

func (s *Server) GetOrder(
	ctx context.Context,
	req *connect.Request[orderv1.GetOrderRequest],
) (*connect.Response[orderv1.GetOrderResponse], error) {
	return s.getOrderHandler.Handle(ctx, req)
}

func (s *Server) ListOrders(
	ctx context.Context,
	req *connect.Request[orderv1.ListOrdersRequest],
) (*connect.Response[orderv1.ListOrdersResponse], error) {
	return s.listOrdersHandler.Handle(ctx, req)
}

func (s *Server) CheckOrderOwner(
	ctx context.Context,
	req *connect.Request[orderv1.CheckOrderOwnerRequest],
) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}
