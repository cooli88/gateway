package orders

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

func TestListOrdersHandler_BackendErrors(t *testing.T) {
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
				listOrdersFunc: func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
					return nil, tt.backendError
				},
			}

			handler := newListOrdersHandler(mockClient)
			_, err := handler.Handle(context.Background(), connect.NewRequest(&orderv1.ListOrdersRequest{}))

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
