package goldmanager

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

func (c Currency) GetField() string {
	r := ""
	switch c {
	case Copper:
		r = "cc"
	case Silver:
		r = "sc"
	case Electrum:
		r = "tc"
	case Gold:
		r = "gc"
	case Platinum:
		r = "pc"
	}
	return r
}

func GetCurrency(s string) Currency {
	var curr Currency
	switch s {
	case "cc":
		curr = Copper
	case "sc":
		curr = Silver
	case "ec":
		curr = Electrum
	case "gc":
		curr = Gold
	case "pc":
		curr = Platinum
	}
	return curr
}

type Coin struct {
	Id       string `json:"_,omitempty"`
	Copper   int
	Silver   int
	Electrum int
	Gold     int
	Platinum int
}

func Convert(value int, src, dest Currency) int {
	if src == dest {
		return value
	}
	div := float64(Converter[src]) / float64(Converter[dest])

	return int(div * float64(value))
}
