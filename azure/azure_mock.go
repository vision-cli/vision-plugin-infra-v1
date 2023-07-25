package azure

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/mock"
)

// TokenCredential
type MockTokenCredential struct {
	mock.Mock
}

// simulates retrieving an access token from the token credential
func (m *MockTokenCredential) GetToken(ctx context.Context, options azcore.TokenCredential) (*azcore.AccessToken, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(*azcore.AccessToken), args.Error(1)
}

// ClientFactory
type MockClientFactory struct {
	mock.Mock
}

// simulates creating a new instance of an armstorage.AccountsClient,
func (m *MockClientFactory) NewAccountsClient() armstorage.AccountsClient {
	args := m.Called()
	return args.Get(0).(armstorage.AccountsClient)
}

// AccountsClient
type MockAccountsClient struct {
	mock.Mock
}

// BeginCreate
func (m *MockAccountsClient) BeginCreate(ctx context.Context, resourceGroupName string, accountName string, accCreateParams armstorage.AccountCreateParameters) (*armstorage.AccountsClientCreateResponse, error) {
	args := m.Called(ctx, resourceGroupName, accountName, accCreateParams)
	// return the expected *armstorage.AccountsClientCreateResponse
	return args.Get(0).(*armstorage.AccountsClientCreateResponse), args.Error(1)
}

// PollingHandler
type MockPollingHandler[T any] struct {
	done   bool
	result T
}

// indicates whether polling operation is complete
func (m *MockPollingHandler[T]) Done() bool {
	return m.done
}

// simulates polling logic
func (m *MockPollingHandler[T]) Poll(ctx context.Context) (*http.Response, error) {
	// simulate the polling logic and return an HTTP response
	return nil, nil
}

// populates out parameter with result value of MockPollingHandler
func (m *MockPollingHandler[T]) Result(ctx context.Context, out *T) error {
	// populate the out parameter with the mock result
	*out = m.result
	return nil
}

// Poller
type MockPoller struct {
	op     MockPollingHandler[armstorage.AccountsClientCreateResponse]
	resp   *http.Response
	err    error
	result armstorage.AccountsClientCreateResponse
	done   bool
}

// simulates waiting for polling operation to complete
func (m *MockPoller) Wait(ctx context.Context) (armstorage.AccountsClientCreateResponse, error) {
	return m.result, nil
}

// returns a pointer to a new MockPoller instance
func NewMockPoller(resp *http.Response, result armstorage.AccountsClientCreateResponse) *MockPoller {
	return &MockPoller{
		op:     MockPollingHandler[armstorage.AccountsClientCreateResponse]{},
		resp:   resp,
		result: result,
	}
}
