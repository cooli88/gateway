package isolation

import (
	"net/http"
	"os"

	"github.com/stretchr/testify/suite"

	"github.com/demo/contracts/gen/go/order/v1/orderv1connect"
)

// Suite is the base suite for isolation tests.
// It provides a Connect RPC client to the Gateway service.
type Suite struct {
	suite.Suite
	orderClient orderv1connect.OrderServiceClient
}

// SetupSuite initializes the Connect RPC client.
func (s *Suite) SetupSuite() {
	gatewayURL := os.Getenv("GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8080"
	}

	s.orderClient = orderv1connect.NewOrderServiceClient(
		http.DefaultClient,
		gatewayURL,
	)
}
