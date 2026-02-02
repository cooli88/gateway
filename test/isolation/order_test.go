package isolation

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/suite"

	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

type OrderTestSuite struct {
	Suite
}

func (s *OrderTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
}

// TestCreateOrder_HappyPath tests successful order creation.
func (s *OrderTestSuite) TestCreateOrder_HappyPath() {
	ctx := context.Background()

	resp, err := s.orderClient.CreateOrder(ctx, connect.NewRequest(&orderv1.CreateOrderRequest{
		UserId: "user-1",
		Item:   "test-item",
		Amount: 100,
	}))

	s.Require().NoError(err)
	s.Equal("order-user-1-001", resp.Msg.Order.Id)
	s.Equal("user-1", resp.Msg.Order.UserId)
	s.Equal("test-item", resp.Msg.Order.Item)
	s.EqualValues(100, resp.Msg.Order.Amount)
}

// TestCreateOrder_ValidationError_EmptyUserId tests validation for empty user_id.
func (s *OrderTestSuite) TestCreateOrder_ValidationError_EmptyUserId() {
	ctx := context.Background()

	_, err := s.orderClient.CreateOrder(ctx, connect.NewRequest(&orderv1.CreateOrderRequest{
		UserId: "",
		Item:   "test-item",
		Amount: 100,
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
}

// TestCreateOrder_ValidationError_EmptyItem tests validation for empty item.
func (s *OrderTestSuite) TestCreateOrder_ValidationError_EmptyItem() {
	ctx := context.Background()

	_, err := s.orderClient.CreateOrder(ctx, connect.NewRequest(&orderv1.CreateOrderRequest{
		UserId: "user-1",
		Item:   "",
		Amount: 100,
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
}

// TestCreateOrder_ValidationError_InvalidAmount tests validation for invalid amount.
func (s *OrderTestSuite) TestCreateOrder_ValidationError_InvalidAmount() {
	ctx := context.Background()

	_, err := s.orderClient.CreateOrder(ctx, connect.NewRequest(&orderv1.CreateOrderRequest{
		UserId: "user-1",
		Item:   "test-item",
		Amount: 0,
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
}

// TestGetOrder_HappyPath tests successful order retrieval.
func (s *OrderTestSuite) TestGetOrder_HappyPath() {
	ctx := context.Background()

	resp, err := s.orderClient.GetOrder(ctx, connect.NewRequest(&orderv1.GetOrderRequest{
		Id:     "order-123",
		UserId: "user-1",
	}))

	s.Require().NoError(err)
	s.Equal("order-123", resp.Msg.Order.Id)
	s.Equal("user-1", resp.Msg.Order.UserId)
}

// TestGetOrder_ValidationError_EmptyId tests validation for empty id.
func (s *OrderTestSuite) TestGetOrder_ValidationError_EmptyId() {
	ctx := context.Background()

	_, err := s.orderClient.GetOrder(ctx, connect.NewRequest(&orderv1.GetOrderRequest{
		Id:     "",
		UserId: "user-1",
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
}

// TestGetOrder_ValidationError_EmptyUserId tests validation for empty user_id.
func (s *OrderTestSuite) TestGetOrder_ValidationError_EmptyUserId() {
	ctx := context.Background()

	_, err := s.orderClient.GetOrder(ctx, connect.NewRequest(&orderv1.GetOrderRequest{
		Id:     "order-123",
		UserId: "",
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodeInvalidArgument, connectErr.Code())
}

// TestGetOrder_NotFound tests that CodeNotFound is returned for non-existent orders.
func (s *OrderTestSuite) TestGetOrder_NotFound() {
	ctx := context.Background()

	_, err := s.orderClient.GetOrder(ctx, connect.NewRequest(&orderv1.GetOrderRequest{
		Id:     "not-found",
		UserId: "user-1",
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodeNotFound, connectErr.Code())
}

// TestGetOrder_PermissionDenied tests that CodePermissionDenied is returned for unauthorized users.
func (s *OrderTestSuite) TestGetOrder_PermissionDenied() {
	ctx := context.Background()

	_, err := s.orderClient.GetOrder(ctx, connect.NewRequest(&orderv1.GetOrderRequest{
		Id:     "order-123",
		UserId: "unauthorized",
	}))

	s.Require().Error(err)
	var connectErr *connect.Error
	s.Require().True(errors.As(err, &connectErr))
	s.Equal(connect.CodePermissionDenied, connectErr.Code())
}

// TestListOrders_HappyPath tests successful orders listing.
func (s *OrderTestSuite) TestListOrders_HappyPath() {
	ctx := context.Background()

	resp, err := s.orderClient.ListOrders(ctx, connect.NewRequest(&orderv1.ListOrdersRequest{}))

	s.Require().NoError(err)
	s.Len(resp.Msg.Orders, 2)
	s.Equal("order-001", resp.Msg.Orders[0].Id)
	s.Equal("order-002", resp.Msg.Orders[1].Id)
}

func TestOrder(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}
