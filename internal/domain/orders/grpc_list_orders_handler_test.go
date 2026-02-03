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

func TestListOrdersHandler(t *testing.T) {
	// Define testData struct locally - Gateway specific
	type testData struct {
		ctx      context.Context
		t        *testing.T
		handler  *listOrdersHandler
		client   *mockOrderServiceClient
		request  *connect.Request[orderv1.ListOrdersRequest]
		response *connect.Response[orderv1.ListOrdersResponse]
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

		// Setup default mock behavior (successful proxy with multiple orders)
		client.listOrdersFunc = func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
			return connect.NewResponse(&orderv1.ListOrdersResponse{
				Orders: []*orderv1.Order{
					{
						Id:     "order-1",
						UserId: "user-123",
						Item:   "Item 1",
						Amount: 100.0,
						Status: "created",
					},
					{
						Id:     "order-2",
						UserId: "user-456",
						Item:   "Item 2",
						Amount: 200.0,
						Status: "completed",
					},
				},
			}), nil
		}

		handler := newListOrdersHandler(client)

		return &testData{
			ctx:     context.Background(),
			t:       t,
			handler: handler,
			client:  client,
		}
	}

	testCases := []testCase{
		// Success scenarios
		{
			name: "Should proxy request to backend successfully",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.ListOrdersRequest{})
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.NoError(td.t, td.err)
				require.NotNil(td.t, td.response)
				require.NotNil(td.t, td.response.Msg.Orders)
				assert.Len(td.t, td.response.Msg.Orders, 2)
				assert.Equal(td.t, "order-1", td.response.Msg.Orders[0].Id)
				assert.Equal(td.t, "order-2", td.response.Msg.Orders[1].Id)
			},
		},
		{
			name: "Should return empty list when backend returns no orders",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.ListOrdersRequest{})
				td.client.listOrdersFunc = func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
					return connect.NewResponse(&orderv1.ListOrdersResponse{
						Orders: []*orderv1.Order{},
					}), nil
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.NoError(td.t, td.err)
				require.NotNil(td.t, td.response)
				assert.Empty(td.t, td.response.Msg.Orders)
			},
		},
		{
			name: "Should return nil orders slice when backend returns nil",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.ListOrdersRequest{})
				td.client.listOrdersFunc = func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
					return connect.NewResponse(&orderv1.ListOrdersResponse{
						Orders: nil,
					}), nil
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.NoError(td.t, td.err)
				require.NotNil(td.t, td.response)
				assert.Nil(td.t, td.response.Msg.Orders)
			},
		},

		// Backend error scenarios
		{
			name: "Should propagate backend Internal error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.ListOrdersRequest{})
				td.client.listOrdersFunc = func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
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
				td.request = connect.NewRequest(&orderv1.ListOrdersRequest{})
				td.client.listOrdersFunc = func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
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
		{
			name: "Should propagate backend DeadlineExceeded error",
			given: func(td *testData) {
				td.request = connect.NewRequest(&orderv1.ListOrdersRequest{})
				td.client.listOrdersFunc = func(ctx context.Context, req *connect.Request[orderv1.ListOrdersRequest]) (*connect.Response[orderv1.ListOrdersResponse], error) {
					return nil, connect.NewError(connect.CodeDeadlineExceeded, errors.New("request timeout"))
				}
			},
			when: func(td *testData) {
				td.response, td.err = td.handler.Handle(td.ctx, td.request)
			},
			then: func(td *testData) {
				require.Error(td.t, td.err)
				assert.Equal(td.t, connect.CodeDeadlineExceeded, connect.CodeOf(td.err))
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
