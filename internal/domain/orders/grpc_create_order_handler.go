package orders

import (
	"context"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
)

type createOrderHandler struct {
	client orderv1connect.OrderServiceClient
}

func newCreateOrderHandler(client orderv1connect.OrderServiceClient) *createOrderHandler {
	return &createOrderHandler{client: client}
}

func (h *createOrderHandler) Handle(
	ctx context.Context,
	req *connect.Request[orderv1.CreateOrderRequest],
) (*connect.Response[orderv1.CreateOrderResponse], error) {
	if err := h.validate(req.Msg); err != nil {
		return nil, err
	}

	return h.client.CreateOrder(ctx, req)
}

func (h *createOrderHandler) validate(req *orderv1.CreateOrderRequest) error {
	if req.UserId == "" {
		return connect.NewError(connect.CodeInvalidArgument, nil)
	}
	if req.Item == "" {
		return connect.NewError(connect.CodeInvalidArgument, nil)
	}
	if req.Amount <= 0 {
		return connect.NewError(connect.CodeInvalidArgument, nil)
	}
	return nil
}
