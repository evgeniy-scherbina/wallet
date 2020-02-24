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

	// check create/get payment(root account -> std account) method
	paymentIDs := test.NewCounterMap()
	for _, accountID := range listAccountIDs {
		// create payment
		id, err := client.CreatePayment(&pbw.CreatePaymentRequest{
			Source:      "root",
			Destination: accountID,
			Amount:      test.DefaultInitialBalance,
		})
		r.Nil(err)
		r.NotEmpty(id)

		paymentIDs.Inc(id)

		// get payment
		payment, err := client.GetPayment(id)
		r.Nil(err)

		r.Equal("root", payment.Source)
		r.Equal(accountID, payment.Destination)
		r.Equal(uint64(test.DefaultInitialBalance), payment.Amount)
	}

	// check get balance method
	for _, accountID := range listAccountIDs {
		balance, err := client.GetBalance(accountID)
		r.Nil(err)

		r.Equal(uint64(test.DefaultInitialBalance), balance)
	}

	// check create/get payment(std account -> std account) method
	paymentNum := 4
	for i := 0; i < paymentNum; i++ {
		// create payment
		source := listAccountIDs[i % 2]
		dest := listAccountIDs[(i + 1) % 2]
		id, err := client.CreatePayment(&pbw.CreatePaymentRequest{
			Source:      source,
			Destination: dest,
			Amount:      test.DefaultPaymentAmount,
		})

		r.Nil(err)
		r.NotEmpty(id)

		paymentIDs.Inc(id)

		// get payment
		payment, err := client.GetPayment(id)
		r.Nil(err)

		r.Equal(source, payment.Source)
		r.Equal(dest, payment.Destination)
		r.Equal(uint64(test.DefaultPaymentAmount), payment.Amount)

		sourceAccountBalance, err := client.GetBalance(source)
		r.Nil(err)
		destAccountBalance, err := client.GetBalance(dest)
		r.Nil(err)

		r.Equal(uint64(test.DefaultPaymentAmount - test.DefaultPaymentAmount), sourceAccountBalance)
		r.Equal(uint64(test.DefaultPaymentAmount + test.DefaultPaymentAmount), destAccountBalance)
	}

	// check list payments method
	{
		payments, err := client.ListPayments()
		r.Nil(err)
		allIDs := test.NewCounterMap()
		for _, payment := range payments {
			allIDs.Inc(payment.Id)
		}

		r.True(allIDs.Contains(paymentIDs))
	}

	// check get balance method(after payments)
	for _, accountID := range listAccountIDs {
		balance, err := client.GetBalance(accountID)
		r.Nil(err)

		r.Equal(int64(test.DefaultInitialBalance), int64(balance))
	}
}