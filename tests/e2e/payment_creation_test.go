package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	wclient "github.com/evgeniy-scherbina/wallet/client"
	"github.com/evgeniy-scherbina/wallet/lib/test"
	pbw "github.com/evgeniy-scherbina/wallet/pb/wallet"
)

func TestPaymentCreation(t *testing.T) {
	r := require.New(t)
	client := wclient.NewHTTPClient(test.DefaultHttpAddress)

	// create root account
	if err := client.CreateRootAccount(); err != nil && !strings.Contains(err.Error(), "can't create root account") {
		r.FailNow(err.Error())
	}

	// check create account method
	accountNum := 2
	accountIDs := test.NewCounterMap()
	for i := 0; i < accountNum; i++ {
		id, err := client.CreateAccount("evgeniy")
		r.Nil(err)
		r.NotEmpty(id)

		accountIDs.Inc(id)
	}
	listAccountIDs := accountIDs.ToSlice()

	// check create payment method
	paymentIDs := test.NewCounterMap()
	for _, accountID := range listAccountIDs {
		id, err := client.CreatePayment(&pbw.CreatePaymentRequest{
			Source:      "root",
			Destination: accountID,
			Amount:      test.DefaultPaymentAmount,
		})
		r.Nil(err)
		r.NotEmpty(id)

		paymentIDs.Inc(id)
	}
	listPaymentIDs := paymentIDs.ToSlice()

	// check get account method
	for _, id := range listPaymentIDs {
		payment, err := client.GetPayment(id)
		r.Nil(err)

		r.Equal("root", payment.Source)
		r.Equal(uint64(test.DefaultPaymentAmount), payment.Amount)
	}
}