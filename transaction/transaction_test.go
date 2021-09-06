package transaction

import (
	"reflect"
	"testing"

	goldmanager "github.com/mjehanno/go-ldenerd-api/gold-manager"
)

type alignTest struct {
	current     goldmanager.Stock
	incoming    goldmanager.Stock
	transaction TransactionType
	expected    goldmanager.Stock
	isFalty     bool
}

var aligns = []alignTest{
	{
		goldmanager.Stock{
			Id:       "",
			Copper:   9,
			Silver:   10,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   10,
			Silver:   0,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		0,
		goldmanager.Stock{
			Id:       "",
			Copper:   9,
			Silver:   9,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		false,
	},
	{
		goldmanager.Stock{
			Id:       "",
			Copper:   10,
			Silver:   0,
			Electrum: 0,
			Gold:     50,
			Platinum: 0,
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   0,
			Silver:   9,
			Electrum: 3,
			Gold:     10,
			Platinum: 0,
		},
		1,
		goldmanager.Stock{
			Id:       "",
			Copper:   10,
			Silver:   9,
			Electrum: 3,
			Gold:     60,
			Platinum: 0,
		},
		false,
	},
	{
		goldmanager.Stock{
			Id:       "",
			Copper:   9,
			Silver:   9,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   9,
			Silver:   9,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		0,
		goldmanager.Stock{
			Id:       "",
			Copper:   0,
			Silver:   0,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		false,
	},
	{
		goldmanager.Stock{
			Id:       "",
			Copper:   9,
			Silver:   0,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   10,
			Silver:   0,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		0,
		goldmanager.Stock{
			Id:       "",
			Copper:   9,
			Silver:   0,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		true,
	},
	{
		goldmanager.Stock{
			Id:       "",
			Copper:   10,
			Silver:   0,
			Electrum: 3,
			Gold:     10,
			Platinum: 0,
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   11,
			Silver:   0,
			Electrum: 0,
			Gold:     0,
			Platinum: 0,
		},
		0,
		goldmanager.Stock{
			Id:       "",
			Copper:   49,
			Silver:   0,
			Electrum: 2,
			Gold:     10,
			Platinum: 0,
		},
		false,
	},
}

/**
* Todo: Corriger cas d'erreur qui doit pas passer dans le deep-equal
 */
func TestAlign(t *testing.T) {
	for _, v := range aligns {
		r, err := Align(v.current, v.incoming, v.transaction)

		if err != nil && !v.isFalty {
			t.Errorf("Failed to align %v et %v on %v type : expected %v. Got %v instead.", v.current, v.incoming, TransactionType(v.transaction), v.expected, err)

			continue
		} else if err == nil && v.isFalty {
			t.Errorf("Failed to align %v et %v on %v type : expected an error. Got %v instead.", v.current, v.incoming, TransactionType(v.transaction), v.expected)
			continue
		} else if !reflect.DeepEqual(r, v.expected) {
			t.Errorf("Failed to align %v et %v on %v type : expected %v, got %v", v.current, v.incoming, v.transaction, v.expected, r)
		}
	}
}

type convertTest struct {
	amount   []Coin
	expected goldmanager.Stock
}

var amounts = []convertTest{
	{
		[]Coin{
			{
				10,
				goldmanager.Copper,
			},
			{
				10,
				goldmanager.Gold,
			},
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   10,
			Silver:   0,
			Electrum: 0,
			Gold:     10,
			Platinum: 0,
		},
	},
	{
		[]Coin{
			{
				5,
				goldmanager.Silver,
			},
			{
				3,
				goldmanager.Platinum,
			},
		},
		goldmanager.Stock{
			Id:       "",
			Copper:   0,
			Silver:   5,
			Electrum: 0,
			Gold:     0,
			Platinum: 3,
		},
	},
}

func TestConvertSumOfAmountToCoin(t *testing.T) {
	for _, v := range amounts {
		if r := ConvertSumOfAmountToCoin(v.amount); !reflect.DeepEqual(r, v.expected) {
			t.Errorf("Conversion of amount %v, to coin failed. Expected : %v. Got : %v", v.amount, v.expected, r)
		}
	}
}
