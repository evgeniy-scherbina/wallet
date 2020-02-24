package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	wclient "github.com/evgeniy-scherbina/wallet/client"
	"github.com/evgeniy-scherbina/wallet/lib/test"
	pbw "github.com/evgeniy-scherbina/wallet/pb/wallet"
)

func TestAmountChecking(t *testing.T) {
	r := require.New(t)
	client := wclient.NewHTTPClient(test.DefaultHttpAddress)

	// create root account
	if err := client.CreateRootAccount(); err != nil && !strings.Contains(err.Error(), "can't create root account") {
		r.FailNow(err.Error())
	}

	// create accounts
	accountNum := 2
	accountIDs := make([]string, 0)
	for i := 0; i < accountNum; i++ {
		id, err := client.CreateAccount("evgeniy")
		r.Nil(err)
		r.NotEmpty(id)

		accountIDs = append(accountIDs, id)
	}
	source := accountIDs[0]
	dest := accountIDs[1]

	// create payment
	id, err := client.CreatePayment(&pbw.CreatePaymentRequest{
		Source:      source,
		Destination: dest,
		Amount:      test.OverflowPaymentAmount,
	})
	r.NotNil(err)
	r.Empty(id)
}