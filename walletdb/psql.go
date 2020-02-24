package walletdb

import (
	"database/sql"
	"encoding/hex"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rogpeppe/fastuuid"

	pbw "github.com/evgeniy-scherbina/wallet/pb/wallet"
)

type DB struct {
	inner *sql.DB
}

func New(user, pass, addr, dbName string) (*DB, error) {
	connString := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", user, pass, addr, dbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	return &DB{
		inner: db,
	}, nil
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (a *Account) ToProto() *pbw.Account {
	return &pbw.Account{
		Id:   a.ID,
		Name: a.Name,
	}
}

type Payment struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Amount      uint64 `json:"amount"`
}

func (p *Payment) ToProto() *pbw.Payment {
	return &pbw.Payment{
		Id:          p.ID,
		Source:      p.Source,
		Destination: p.Destination,
		Amount:      p.Amount,
	}
}

func (db *DB) PutAccount(account *Account) (string, error) {
	if account.ID == "" || len(account.ID) == 0 {
		var id [24]byte
		id = fastuuid.MustNewGenerator().Next()
		account.ID = hex.EncodeToString(id[:])
	}

	_, err := db.inner.Exec("INSERT INTO accounts VALUES($1, $2)", account.ID, account.Name)
	if err != nil {
		return "", err
	}

	return account.ID, nil
}

func (db *DB) GetAccount(id string) (*Account, error) {
	row := db.inner.QueryRow("SELECT * FROM accounts WHERE id = $1", id)

	account := new(Account)
	if err := row.Scan(&account.ID, &account.Name); err != nil {
		return nil, err
	}

	return account, nil
}

func (db *DB) ListAccounts() ([]*Account, error) {
	rows, err := db.inner.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*Account, 0)
	for rows.Next() {
		account := new(Account)
		if err := rows.Scan(&account.ID, &account.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (db *DB) PutPayment(payment *Payment) (string, error) {
	if payment.ID == "" || len(payment.ID) == 0 {
		var id [24]byte
		id = fastuuid.MustNewGenerator().Next()
		payment.ID = hex.EncodeToString(id[:])
	}

	_, err := db.inner.Exec(
		"INSERT INTO payments VALUES($1, $2, $3, $4)",
		payment.ID,
		payment.Source,
		payment.Destination,
		payment.Amount,
	)
	if err != nil {
		return "", err
	}

	return payment.ID, nil
}

func (db *DB) GetPayment(id string) (*Payment, error) {
	row := db.inner.QueryRow("SELECT * FROM payments WHERE id = $1", id)

	payment := new(Payment)
	if err := row.Scan(&payment.ID, &payment.Source, &payment.Destination, &payment.Amount); err != nil {
		return nil, err
	}

	return payment, nil
}

func (db *DB) ListPayments() ([]*Payment, error) {
	rows, err := db.inner.Query("SELECT * FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]*Payment, 0)
	for rows.Next() {
		payment := new(Payment)
		if err := rows.Scan(&payment.ID, &payment.Source, &payment.Destination, &payment.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (db *DB) GetBalance(accountId string) (uint64, error) {
	deposit, err := db.deposit(accountId)
	if err != nil {
		return 0, err
	}

	credit, err := db.credit(accountId)
	if err != nil {
		return 0, err
	}
	_ = credit

	return deposit, nil
}

func (db *DB) deposit(accountId string) (uint64, error) {
	row := db.inner.QueryRow("SELECT SUM(amount) FROM payments WHERE destination = $1", accountId)

	var sum sql.NullInt64
	if err := row.Scan(&sum); err != nil {
		return 0, err
	}

	if sum.Valid {
		return uint64(sum.Int64), nil
	} else {
		return 0, nil
	}
}

func (db *DB) credit(accountId string) (uint64, error) {
	row := db.inner.QueryRow("SELECT SUM(amount) FROM payments WHERE source = $1", accountId)

	var sum sql.NullInt64
	if err := row.Scan(&sum); err != nil {
		return 0, err
	}

	if sum.Valid {
		return uint64(sum.Int64), nil
	} else {
		return 0, nil
	}
}