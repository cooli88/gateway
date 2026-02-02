package orders

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

func TestCreateOrderHandler_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request *orderv1.CreateOrderRequest
	}{
		{
			name: "validation_error_empty_user_id",
			request: &orderv1.CreateOrderRequest{
				UserId: "",
				Item:   "test-item",
				Amount: 100,
			},
		},
		{
			name: "validation_error_empty_item",
			request: &orderv1.CreateOrderRequest{
				UserId: "user-123",
				Item:   "",
				Amount: 100,
			},
		},
		{
			name: "validation_error_zero_amount",
			request: &orderv1.CreateOrderRequest{
				UserId: "user-123",
				Item:   "test-item",
				Amount: 0,
			},
		},
		{
			name: "validation_error_negative_amount",
			request: &orderv1.CreateOrderRequest{
				UserId: "user-123",
				Item:   "test-item",
				Amount: -100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backendCalled := false
			mockClient := &mockOrderServiceClient{
				createOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
					backendCalled = true
					return connect.NewResponse(&orderv1.CreateOrderResponse{}), nil
				},
			}

			handler := newCreateOrderHandler(mockClient)
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

			if backendCalled {
				t.Error("backend should not be called on validation error")
			}
		})
	}
}

func TestCreateOrderHandler_BackendErrors(t *testing.T) {
	tests := []struct {
		name         string
		backendError *connect.Error
		expectedCode connect.Code
	}{
		{
			name:         "backend_error_internal",
			backendError: connect.NewError(connect.CodeInternal, nil),
			expectedCode: connect.CodeInternal,
		},
		{
			name:         "backend_error_unavailable",
			backendError: connect.NewError(connect.CodeUnavailable, nil),
			expectedCode: connect.CodeUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockOrderServiceClient{
				createOrderFunc: func(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
					return nil, tt.backendError
				},
			}

			handler := newCreateOrderHandler(mockClient)
			_, err := handler.Handle(context.Background(), connect.NewRequest(&orderv1.CreateOrderRequest{
				UserId: "user-123",
				Item:   "test-item",
				Amount: 100,
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
		})
	}
}
