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

func TestCreateOrderHandler(t *testing.T) {
	// Define testData struct locally - Gateway specific
	type testData struct {
		ctx      context.Context
		t        *testing.T
		handler  *createOrderHandler
		client   *mockOrderServiceClient
		request  *connect.Request[orderv1.CreateOrderRequest]
		response *connect.Response[orderv1.CreateOrderResponse]
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

		// Setup default mock behavior (successful proxy)
		client.createOrderFunc = func(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
			return connect.NewResponse(&orderv1.CreateOrderResponse{
				Order: &orderv1.Order{
					Id:     "order-123",
					UserId: req.Msg.UserId,
					Item:   req.Msg.Item,
					Amount: req.Msg.Amount,
					Status: "created",
				},
			}), nil
		}

		handler := newCreateOrderHandler(client)

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
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "user-123",
					Item:   "Test Item",
					Amount: 100.50,
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
				assert.Equal(td.t, "user-123", td.response.Msg.Order.UserId)
				assert.Equal(td.t, "Test Item", td.response.Msg.Order.Item)
				assert.Equal(td.t, 100.50, td.response.Msg.Order.Amount)
			},
		},

		// Validation errors
		{
			name: "Should return InvalidArgument when user_id is empty",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "",
					Item:   "Test Item",
					Amount: 100,
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
			name: "Should return InvalidArgument when item is empty",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "user-123",
					Item:   "",
					Amount: 100,
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
			name: "Should return InvalidArgument when amount is zero",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "user-123",
					Item:   "Test Item",
					Amount: 0,
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
			name: "Should return InvalidArgument when amount is negative",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "user-123",
					Item:   "Test Item",
					Amount: -50,
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

		// Backend error scenarios
		{
			name: "Should propagate backend Internal error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "user-123",
					Item:   "Test Item",
					Amount: 100,
				})
				td.client.createOrderFunc = func(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
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
		{
			name: "Should propagate backend Unavailable error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.CreateOrderRequest{
					UserId: "user-123",
					Item:   "Test Item",
					Amount: 100,
				})
				td.client.createOrderFunc = func(ctx context.Context, req *connect.Request[orderv1.CreateOrderRequest]) (*connect.Response[orderv1.CreateOrderResponse], error) {
					return nil, connect.NewError(connect.CodeUnavailable, errors.New("service unavailable"))
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeUnavailable, connect.CodeOf(td.err))
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
