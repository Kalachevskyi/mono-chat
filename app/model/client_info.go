package model

// ClientInfo - represents clients information struct.
type ClientInfo struct {
	ClientID string `json:"clientId"`
	Name     string `json:"name"`
	Accounts []struct {
		ID           string   `json:"id"`
		CurrencyCode int      `json:"currencyCode"`
		CashbackType string   `json:"cashbackType"`
		Balance      int      `json:"balance"`
		CreditLimit  int      `json:"creditLimit"`
		MaskedPan    []string `json:"maskedPan"`
		Type         string   `json:"type"`
	} `json:"accounts"`
}
