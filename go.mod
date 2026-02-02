module github.com/demo/gateway

go 1.21

require (
	connectrpc.com/connect v1.16.2
	github.com/demo/contracts v0.0.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/net v0.23.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/demo/contracts => ../contracts
