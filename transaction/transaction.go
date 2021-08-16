package transaction

type Transaction struct {
	Type   TransactionType
	Value  float64
	Reason string
}

type TransactionType int

const (
	Debit TransactionType = iota
	Credit
)
