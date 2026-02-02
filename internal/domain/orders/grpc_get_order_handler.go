package orders

import (
	"context"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
)

type getOrderHandler struct {
	client orderv1connect.OrderServiceClient
}

func newGetOrderHandler(client orderv1connect.OrderServiceClient) *getOrderHandler {
	return &getOrderHandler{client: client}
}

func (h *getOrderHandler) Handle(
	ctx context.Context,
	req *connect.Request[orderv1.GetOrderRequest],
) (*connect.Response[orderv1.GetOrderResponse], error) {
	if err := h.validate(req.Msg); err != nil {
		return nil, err
	}

	_, err := h.client.CheckOrderOwner(ctx, connect.NewRequest(&orderv1.CheckOrderOwnerRequest{
		OrderId: req.Msg.Id,
		UserId:  req.Msg.UserId,
	}))
	if err != nil {
		return nil, err
	}

	return h.client.GetOrder(ctx, req)
}

func (h *getOrderHandler) validate(req *orderv1.GetOrderRequest) error {
	if req.Id == "" {
		return connect.NewError(connect.CodeInvalidArgument, nil)
	}
	if req.UserId == "" {
		return connect.NewError(connect.CodeInvalidArgument, nil)
	}
	return nil
}
