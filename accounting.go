package accounting

import (
	"time"
)

type Order struct {
	Date   time.Time
	Symbol string
	Amount int
	Fees   int
	Pieces float64
}

type Transfer struct {
	Date        time.Time
	Source      string
	Target      string
	Description string
	Amount      int
}

type Account struct {
	IBAN string
	Name string
}
