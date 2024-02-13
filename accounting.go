package accounting

import (
	"os"
	"path/filepath"
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

func ListTxts(rootpath string) ([]string, error) {
	list := []string{}

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".txt" {
			list = append(list, path)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return list, nil
}
