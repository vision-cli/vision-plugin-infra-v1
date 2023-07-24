package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (m *MockAccountsClient) BeginCreate(ctx context.Context, resourceGroupName string, accountName string, accCreateParams armstorage.AccountCreateParameters) (*armstorage.AccountsClientCreateResponse, error) {
	args := m.Called(ctx, resourceGroupName, accountName, accCreateParams)
	// )return the expected *armstorage.AccountsClientCreateResponse
	return args.Get(0).(*armstorage.AccountsClientCreateResponse), args.Error(1)
}

func TestBeginCreate(t *testing.T) {
	var mac MockAccountsClient
	mockResult := &armstorage.AccountsClientCreateResponse{}

	mac.On("BeginCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(mockResult, nil)

	result, err := mac.BeginCreate(ctx, resourceGroupName, accountName, accCreateParams)
	require.NoError(t, err)

	assert.Equal(t, mockResult, result)

	// check if appropriate methods were called
	mac.AssertExpectations(t)
}
