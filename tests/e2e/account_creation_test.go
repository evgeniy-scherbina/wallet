package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	wclient "github.com/evgeniy-scherbina/wallet/client"
	"github.com/evgeniy-scherbina/wallet/lib/test"
)

func TestAccountCreation(t *testing.T) {
	r := require.New(t)
	client := wclient.NewHTTPClient(test.DefaultHttpAddress)

	// check create account method
	accountNum := 2
	ids := test.NewCounterMap()
	for i := 0; i < accountNum; i++ {
		id, err := client.CreateAccount("evgeniy")
		r.Nil(err)
		r.NotEmpty(id)

		ids.Inc(id)
	}

	// check get account method
	for key := range ids.ToMap() {
		account, err := client.GetAccount(key)
		r.Nil(err)
		r.Equal("evgeniy", account.Name)
	}

	// check list account methods
	accounts, err := client.ListAccounts()
	r.Nil(err)
	allIDs := test.NewCounterMap()
	for _, account := range accounts {
		allIDs.Inc(account.Id)
	}

	r.True(allIDs.Contains(ids))
}