package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	pbw "github.com/evgeniy-scherbina/wallet/pb/wallet"
)

const (
	createAccountEndpoint     = "/api/v1/account"
	createRootAccountEndpoint = "/api/v1/account/root"
	getAccountEndpoint        = "/api/v1/account"
	listAccountsEndpoint      = "/api/v1/list/account"

	createPaymentEndpoint = "/api/v1/payment"
	getPaymentEndpoint    = "/api/v1/payment"
	listPaymentEndpoint   = "/api/v1/list/payment"
)

type Payment struct {
	Id          string `json:"id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Amount      uint64 `json:"amount,string"`
}

type ListPaymentsResponse struct {
	Payments []*Payment `json:"payments"`
}

type httpClient struct {
	httpAddress string
}

func NewHTTPClient(httpAddress string) *httpClient {
	return &httpClient{
		httpAddress: httpAddress,
	}
}

func (httpClient *httpClient) CreateAccount(name string) (string, error) {
	url := httpClient.httpAddress + createAccountEndpoint

	req := pbw.CreateAccountRequest{
		Name: name,
	}
	raw, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBuffer(raw)

	httpReq, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return "", err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}
	if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	var resp pbw.CreateAccountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (httpClient *httpClient) CreateRootAccount() error {
	url := httpClient.httpAddress + createRootAccountEndpoint

	buff := bytes.NewBufferString("{}")
	httpReq, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 200 {
		return fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	return nil
}

func (httpClient *httpClient) GetAccount(id string) (*pbw.Account, error) {
	url := fmt.Sprintf("%v%v/%v", httpClient.httpAddress, getAccountEndpoint, id)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	var resp pbw.Account
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (httpClient *httpClient) ListAccounts() ([]*pbw.Account, error) {
	url := httpClient.httpAddress + listAccountsEndpoint

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	var resp pbw.ListAccountsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Accounts, nil
}

func (httpClient *httpClient) CreatePayment(req *pbw.CreatePaymentRequest) (string, error) {
	url := httpClient.httpAddress + createPaymentEndpoint

	raw, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBuffer(raw)

	httpReq, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return "", err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}
	if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	var resp pbw.CreateAccountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (httpClient *httpClient) GetPayment(id string) (*Payment, error) {
	url := fmt.Sprintf("%v%v/%v", httpClient.httpAddress, getPaymentEndpoint, id)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	var resp Payment
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (httpClient *httpClient) ListPayments() ([]*Payment, error) {
	url := httpClient.httpAddress + listPaymentEndpoint

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("wrong status code, want 200, got %v, details: %s", httpResp.StatusCode, body)
	}

	var resp ListPaymentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Payments, nil
}