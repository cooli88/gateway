package orders

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

func TestGetOrderHandler_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request *orderv1.GetOrderRequest
	}{
		{
			name: "validation_error_empty_id",
			request: &orderv1.GetOrderRequest{
				Id:     "",
				UserId: "user-123",
			},
		},
		{
			name: "validation_error_empty_user_id",
			request: &orderv1.GetOrderRequest{
				Id:     "order-123",
				UserId: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkOwnerCalled := false
			getOrderCalled := false

			mockClient := &mockOrderServiceClient{
				checkOrderOwnerFunc: func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
					checkOwnerCalled = true
					return connect.NewResponse(&orderv1.CheckOrderOwnerResponse{}), nil
				},
				getOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
					getOrderCalled = true
					return connect.NewResponse(&orderv1.GetOrderResponse{}), nil
				},
			}

			handler := newGetOrderHandler(mockClient)
			_, err := handler.Handle(context.Background(), connect.NewRequest(tt.request))

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			connectErr, ok := err.(*connect.Error)
			if !ok {
				t.Fatalf("expected connect.Error, got %T", err)
			}

			if connectErr.Code() != connect.CodeInvalidArgument {
				t.Errorf("expected InvalidArgument, got %v", connectErr.Code())
			}

			if checkOwnerCalled {
				t.Error("CheckOrderOwner should not be called on validation error")
			}

			if getOrderCalled {
				t.Error("GetOrder should not be called on validation error")
			}
		})
	}
}

func TestGetOrderHandler_CheckOwnerErrors(t *testing.T) {
	tests := []struct {
		name         string
		ownerError   *connect.Error
		expectedCode connect.Code
	}{
		{
			name:         "check_owner_permission_denied",
			ownerError:   connect.NewError(connect.CodePermissionDenied, nil),
			expectedCode: connect.CodePermissionDenied,
		},
		{
			name:         "check_owner_not_found",
			ownerError:   connect.NewError(connect.CodeNotFound, nil),
			expectedCode: connect.CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getOrderCalled := false

			mockClient := &mockOrderServiceClient{
				checkOrderOwnerFunc: func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
					return nil, tt.ownerError
				},
				getOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
					getOrderCalled = true
					return connect.NewResponse(&orderv1.GetOrderResponse{}), nil
				},
			}

			handler := newGetOrderHandler(mockClient)
			_, err := handler.Handle(context.Background(), connect.NewRequest(&orderv1.GetOrderRequest{
				Id:     "order-123",
				UserId: "user-123",
			}))

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			connectErr, ok := err.(*connect.Error)
			if !ok {
				t.Fatalf("expected connect.Error, got %T", err)
			}

			if connectErr.Code() != tt.expectedCode {
				t.Errorf("expected %v, got %v", tt.expectedCode, connectErr.Code())
			}

			if getOrderCalled {
				t.Error("GetOrder should not be called when CheckOrderOwner fails")
			}
		})
	}
}

func TestGetOrderHandler_GetOrderError(t *testing.T) {
	t.Run("get_order_internal_error", func(t *testing.T) {
		checkOwnerCalled := false

		mockClient := &mockOrderServiceClient{
			checkOrderOwnerFunc: func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
				checkOwnerCalled = true
				return connect.NewResponse(&orderv1.CheckOrderOwnerResponse{}), nil
			},
			getOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
				return nil, connect.NewError(connect.CodeInternal, nil)
			},
		}

		handler := newGetOrderHandler(mockClient)
		_, err := handler.Handle(context.Background(), connect.NewRequest(&orderv1.GetOrderRequest{
			Id:     "order-123",
			UserId: "user-123",
		}))

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !checkOwnerCalled {
			t.Error("CheckOrderOwner should be called before GetOrder")
		}

		connectErr, ok := err.(*connect.Error)
		if !ok {
			t.Fatalf("expected connect.Error, got %T", err)
		}

		if connectErr.Code() != connect.CodeInternal {
			t.Errorf("expected Internal, got %v", connectErr.Code())
		}
	})
}

func TestGetOrderHandler_CallOrder(t *testing.T) {
	t.Run("call_order_verification", func(t *testing.T) {
		var callSequence []string

		mockClient := &mockOrderServiceClient{
			checkOrderOwnerFunc: func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
				callSequence = append(callSequence, "CheckOrderOwner")
				return connect.NewResponse(&orderv1.CheckOrderOwnerResponse{}), nil
			},
			getOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
				callSequence = append(callSequence, "GetOrder")
				return connect.NewResponse(&orderv1.GetOrderResponse{
					Order: &orderv1.Order{
						Id:     "order-123",
						UserId: "user-123",
						Item:   "test-item",
						Amount: 100,
					},
				}), nil
			},
		}

		handler := newGetOrderHandler(mockClient)
		_, err := handler.Handle(context.Background(), connect.NewRequest(&orderv1.GetOrderRequest{
			Id:     "order-123",
			UserId: "user-123",
		}))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(callSequence) != 2 {
			t.Fatalf("expected 2 calls, got %d", len(callSequence))
		}

		if callSequence[0] != "CheckOrderOwner" {
			t.Errorf("expected first call to be CheckOrderOwner, got %s", callSequence[0])
		}

		if callSequence[1] != "GetOrder" {
			t.Errorf("expected second call to be GetOrder, got %s", callSequence[1])
		}
	})
}
