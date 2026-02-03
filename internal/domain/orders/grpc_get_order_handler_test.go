package orders

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderv1 "github.com/demo/contracts/gen/go/order/v1"
)

func TestGetOrderHandler(t *testing.T) {
	// Define testData struct locally - Gateway specific
	type testData struct {
		ctx      context.Context
		t        *testing.T
		handler  *getOrderHandler
		client   *mockOrderServiceClient
		request  *connect.Request[orderv1.GetOrderRequest]
		response *connect.Response[orderv1.GetOrderResponse]
		err      error
	}

	// Define testCase struct locally - GWT pattern is MANDATORY
	type testCase struct {
		name  string
		given func(*testData)
		when  func(*testData)
		then  func(*testData)
	}

	// Setup function creates isolated test data for each test case
	setupTestData := func(t *testing.T) *testData {
		client := &mockOrderServiceClient{}

		// Setup default mock behavior for CheckOrderOwner (success)
		client.checkOrderOwnerFunc = func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
			return connect.NewResponse(&orderv1.CheckOrderOwnerResponse{}), nil
		}

		// Setup default mock behavior for GetOrder (success)
		client.getOrderFunc = func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
			return connect.NewResponse(&orderv1.GetOrderResponse{
				Order: &orderv1.Order{
					Id:     req.Msg.Id,
					UserId: req.Msg.UserId,
					Item:   "Test Item",
					Amount: 100.50,
					Status: "created",
				},
			}), nil
		}

		handler := newGetOrderHandler(client)

		return &testData{
			ctx:     context.Background(),
			t:       t,
			handler: handler,
			client:  client,
		}
	}

	testCases := []testCase{
		// Success scenario
		{
			name: "Should proxy request to backend successfully",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "order-123",
					UserId: "user-456",
				})
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.NoError(td.t, td.err)
				require.NotNil(td.t, td.response)
				require.NotNil(td.t, td.response.Msg.Order)
				assert.Equal(td.t, "order-123", td.response.Msg.Order.Id)
				assert.Equal(td.t, "user-456", td.response.Msg.Order.UserId)
				assert.Equal(td.t, "Test Item", td.response.Msg.Order.Item)
			},
		},

		// Validation errors
		{
			name: "Should return InvalidArgument when id is empty",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "",
					UserId: "user-456",
				})
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeInvalidArgument, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},
		{
			name: "Should return InvalidArgument when user_id is empty",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "order-123",
					UserId: "",
				})
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeInvalidArgument, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},
		{
			name: "Should return InvalidArgument when both id and user_id are empty",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "",
					UserId: "",
				})
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeInvalidArgument, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},

		// CheckOrderOwner error scenarios
		{
			name: "Should propagate CheckOrderOwner NotFound error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "order-123",
					UserId: "user-456",
				})
				td.client.checkOrderOwnerFunc = func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
					return nil, connect.NewError(connect.CodeNotFound, errors.New("order not found"))
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeNotFound, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},
		{
			name: "Should propagate CheckOrderOwner PermissionDenied error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "order-123",
					UserId: "user-456",
				})
				td.client.checkOrderOwnerFunc = func(ctx context.Context, req *connect.Request[orderv1.CheckOrderOwnerRequest]) (*connect.Response[orderv1.CheckOrderOwnerResponse], error) {
					return nil, connect.NewError(connect.CodePermissionDenied, errors.New("user is not the owner"))
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodePermissionDenied, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},

		// GetOrder backend error scenarios
		{
			name: "Should propagate GetOrder NotFound error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "order-123",
					UserId: "user-456",
				})
				td.client.getOrderFunc = func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
					return nil, connect.NewError(connect.CodeNotFound, errors.New("order not found"))
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeNotFound, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},
		{
			name: "Should propagate GetOrder Internal error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.GetOrderRequest{
					Id:     "order-123",
					UserId: "user-456",
				})
				td.client.getOrderFunc = func(ctx context.Context, req *connect.Request[orderv1.GetOrderRequest]) (*connect.Response[orderv1.GetOrderResponse], error) {
					return nil, connect.NewError(connect.CodeInternal, errors.New("database error"))
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeInternal, connect.CodeOf(td.err))
				assert.Nil(td.t, td.response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			td := setupTestData(t)
			td.t = t
			tc.given(td)
			tc.when(td)
			tc.then(td)
		})
	}
}
