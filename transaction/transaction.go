package transaction

import goldmanager "github.com/mjehanno/goldenerd/gold-manager"

type Transaction struct {
	Id     string `json:"_id,omitempty"`
	Type   TransactionType
	Amount []Amount
	IsGem  bool
	Reason string
}

type TransactionType int

type Amount struct {
	Value    int
	Currency goldmanager.Currency
}

const (
	Debit TransactionType = iota
	Credit
)
