package goldmanager

//Represent a currency, those are based on D&d 5e rules.
type Currency int

const (
	Copper Currency = iota
	Silver
	Electrum
	Gold
	Platinum
	Limit
)

var Converter = map[Currency]int{
	Copper:   1,
	Silver:   10,
	Electrum: 50,
	Gold:     100,
	Platinum: 1000,
}

func (c Currency) String() string {
	r := ""
	switch c {
	case Copper:
		r = "Copper"
	case Silver:
		r = "Silver"
	case Electrum:
		r = "Electrum"
	case Gold:
		r = "Gold"
	case Platinum:
		r = "Platinum"
	}
	return r
}

//Represent a set of coins
type Stock struct {
	Id       string `json:"_,omitempty"`
	Copper   int
	Silver   int
	Electrum int
	Gold     int
	Platinum int
}

//Convert coins
func Convert(value int, src, dest Currency) int {
	if src == dest {
		return value
	}
	div := float64(Converter[src]) / float64(Converter[dest])

	return int(div * float64(value))
}
