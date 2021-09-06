package transaction

import (
	"errors"
	"reflect"

	goldmanager "github.com/mjehanno/go-ldenerd-api/gold-manager"
)

//Represent a transaction made by any player or group of player
type Transaction struct {
	Id     string `json:"_id,omitempty"`
	Type   TransactionType
	Amount []Coin
	IsGem  bool
	Reason string
}

//Represent a type of transaction
type TransactionType int

//Represent an amount in a transaction
type Coin struct {
	Value    int
	Currency goldmanager.Currency
}

const (
	Debit TransactionType = iota
	Credit
)

func (t TransactionType) String() string {
	r := ""
	switch t {
	case Debit:
		r = "Debit"
	case Credit:
		r = "Credit"
	}
	return r
}

//Convert a list of amount coming from a transaction
func ConvertSumOfAmountToCoin(amounts []Coin) goldmanager.Stock {
	var c goldmanager.Stock

	reflectValue := reflect.ValueOf(&c)
	coinValue := reflectValue.Elem()
	for _, a := range amounts {
		currentCoin := coinValue.FieldByName(a.Currency.String())
		sum := currentCoin.Int() + int64(a.Value)
		currentCoin.SetInt(sum)
	}
	return c
}

//Do the math on transaction
func Align(current, incoming goldmanager.Stock, tr TransactionType) (goldmanager.Stock, error) {
	if tr == Credit {
		current.Copper += incoming.Copper
		current.Silver += incoming.Silver
		current.Electrum += incoming.Electrum
		current.Gold += incoming.Gold
		current.Platinum += incoming.Platinum
	} else {
		incomingAmounts := map[int]reflect.Value{}
		inc := reflect.ValueOf(incoming)

		currentAmount := reflect.ValueOf(&current)
		curr := currentAmount.Elem()

		for i := 1; i < inc.NumField(); i++ {
			if inc.Field(i).Int() > 0 {
				incomingAmounts[i] = inc.Field(i)
			}
		}

		for currency, value := range incomingAmounts {
			if curr.Field(currency).Int() >= value.Int() {
				curr.Field(currency).SetInt(curr.Field(currency).Int() - value.Int())
			} else {
				for i := currency + 1; i < curr.NumField(); i++ {
					if i == curr.NumField()-1 {
						return current, errors.New("cot enough coins")
					}
					if curr.Field(i).Int() > 0 {
						curr.Field(i).SetInt(curr.Field(i).Int() - 1)
						newVal := curr.Field(currency).Int() + int64(goldmanager.Convert(1, goldmanager.Currency(i-1), goldmanager.Currency(currency-1))) - value.Int()
						curr.Field(currency).SetInt(newVal)
						break
					} else {
						continue
					}

				}
			}
		}
	}
	return current, nil
}
