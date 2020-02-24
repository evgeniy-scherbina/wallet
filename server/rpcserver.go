package server

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"

	pbw "github.com/evgeniy-scherbina/wallet/pb/wallet"
	"github.com/evgeniy-scherbina/wallet/walletdb"
)

type WalletServer struct {
	db *walletdb.DB
}

func NewWalletServer(db *walletdb.DB) *WalletServer {
	return &WalletServer{
		db: db,
	}
}

func (ws *WalletServer) CreateAccount(ctx context.Context, req *pbw.CreateAccountRequest) (*pbw.CreateAccountResponse, error) {
	id, err := ws.db.PutAccount(&walletdb.Account{
		Name: req.Name,
	})
	if err != nil {
		err = fmt.Errorf("can't create account: %v", err)
		log.Error(err)
		return nil, err
	}

	return &pbw.CreateAccountResponse{
		Id: id,
	}, nil
}

func (ws *WalletServer) CreateRootAccount(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	_, err := ws.db.PutAccount(&walletdb.Account{
		ID:   "root",
		Name: "root",
	})
	if err != nil {
		err = fmt.Errorf("can't create root account: %v", err)
		log.Error(err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (ws *WalletServer) GetAccount(ctx context.Context, req *pbw.GetAccountRequest) (*pbw.Account, error) {
	account, err := ws.db.GetAccount(req.Id)
	if err != nil {
		err = fmt.Errorf("can't get account: %v", err)
		log.Error(err)
		return nil, err
	}

	return account.ToProto(), nil
}

func (ws *WalletServer) ListAccounts(ctx context.Context, _ *empty.Empty) (*pbw.ListAccountsResponse, error) {
	accounts, err := ws.db.ListAccounts()
	if err != nil {
		err = fmt.Errorf("can't list accounts: %v", err)
		log.Error(err)
		return nil, err
	}

	resp := new(pbw.ListAccountsResponse)
	for _, account := range accounts {
		resp.Accounts = append(resp.Accounts, account.ToProto())
	}

	return resp, nil
}

func (ws *WalletServer) CreatePayment(ctx context.Context, req *pbw.CreatePaymentRequest) (*pbw.CreatePaymentResponse, error) {
	balance, err := ws.db.GetBalance(req.Source)
	if err != nil {
		err = fmt.Errorf("can't get balance: %v", err)
		log.Error(err)
		return nil, err
	}
	if balance < req.Amount {
		log.Warn(ErrInsufficientFunds)
		return nil, ErrInsufficientFunds
	}

	id, err := ws.db.PutPayment(&walletdb.Payment{
		Source:      req.Source,
		Destination: req.Destination,
		Amount:      req.Amount,
	})
	if err != nil {
		err = fmt.Errorf("can't create payment: %v", err)
		log.Error(err)
		return nil, err
	}

	return &pbw.CreatePaymentResponse{
		Id: id,
	}, nil
}

func (ws *WalletServer) GetPayment(ctx context.Context, req *pbw.GetPaymentRequest) (*pbw.Payment, error) {
	payment, err := ws.db.GetPayment(req.Id)
	if err != nil {
		err = fmt.Errorf("can't get payment: %v", err)
		log.Error(err)
		return nil, err
	}

	return payment.ToProto(), nil
}

func (ws *WalletServer) ListPayments(ctx context.Context, _ *empty.Empty) (*pbw.ListPaymentsResponse, error) {
	payments, err := ws.db.ListPayments()
	if err != nil {
		err = fmt.Errorf("can't list payments: %v", err)
		log.Error(err)
		return nil, err
	}

	resp := new(pbw.ListPaymentsResponse)
	for _, payment := range payments {
		resp.Payments = append(resp.Payments, payment.ToProto())
	}

	return resp, nil
}

func (ws *WalletServer) GetBalance(ctx context.Context, req *pbw.GetBalanceRequest) (*pbw.GetBalanceResponse, error) {
	balance, err := ws.db.GetBalance(req.AccountId)
	if err != nil {
		err = fmt.Errorf("can't get balance: %v", err)
		log.Error(err)
		return nil, err
	}

	return &pbw.GetBalanceResponse{
		Amount: balance,
	}, nil
}