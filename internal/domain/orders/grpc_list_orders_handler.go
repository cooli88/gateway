package orders

import (
	"context"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
)

type listOrdersHandler struct {
	client orderv1connect.OrderServiceClient
}

func newListOrdersHandler(client orderv1connect.OrderServiceClient) *listOrdersHandler {
	return &listOrdersHandler{client: client}
}

func (h *listOrdersHandler) Handle(
	ctx context.Context,
	req *connect.Request[orderv1.ListOrdersRequest],
) (*connect.Response[orderv1.ListOrdersResponse], error) {
	if err := h.validate(req.Msg); err != nil {
		return nil, err
	}

	return h.client.ListOrders(ctx, req)
}

func (h *listOrdersHandler) validate(req *orderv1.ListOrdersRequest) error {
	return nil
}
