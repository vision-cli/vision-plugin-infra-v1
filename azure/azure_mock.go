package azure_test

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/mock"
)

// mock struct that implements AccountsClient interface
type MockAccountsClient struct {
	mock.Mock
}

// mock struct that implements Pollinghandler
type MockPollingHandler[T any] struct {
	done   bool
	result T
}

// mock struct that implements Poller interface
type MockPoller struct {
	op     MockPollingHandler[armstorage.AccountsClientCreateResponse]
	resp   *http.Response
	err    error
	result armstorage.AccountsClientCreateResponse
	done   bool
}

func (m *MockPollingHandler[T]) Done() bool {
	return m.done
}

func (m *MockPollingHandler[T]) Poll(ctx context.Context) (*http.Response, error) {
	// simulate the polling logic and return an HTTP response
	return nil, nil
}

func (m *MockPollingHandler[T]) Result(ctx context.Context, out *T) error {
	// populate the out parameter with the mock result
	*out = m.result
	return nil
}

func (m *MockPoller) Wait(ctx context.Context) (armstorage.AccountsClientCreateResponse, error) {
	return m.result, nil
}

// create the mock poller
func NewMockPoller(resp *http.Response, result armstorage.AccountsClientCreateResponse) *MockPoller {
	return &MockPoller{
		op:     MockPollingHandler[armstorage.AccountsClientCreateResponse]{},
		resp:   resp,
		result: result,
	}
}
