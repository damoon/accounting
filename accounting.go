package accounting

import (
	"time"
)

type Order struct {
	Date   time.Time
	Symbol string
	Price  int
	Fees   int
	Pieces float64
}

type Transfer struct {
	Date        time.Time
	Description string
	Amount      int
}

type Account struct {
	IBAN      string
	Name      string
	Transfers []Transfer
}
