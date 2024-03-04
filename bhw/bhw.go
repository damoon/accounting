package bhw

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/damoon/accounting"
)

type BHW struct {
	accounts map[string]accounting.Account
}

func NewBhw() BHW {
	return BHW{
		accounts: map[string]accounting.Account{},
	}
}

func (bhw *BHW) LoadFrom(p accounting.Pdf) error {
	content, err := p.WithLayout()
	if err != nil {
		return err
	}

	if strings.Contains(content, "Wohn Plus") || strings.Contains(content, "BHW WohnBausparen Plus") {
		return bhw.load(content)
	}

	return nil
}

var transfersRegex = regexp.MustCompile(`(\d\d\.\d\d\.\d\d)  \s+  (\S+(\s?\S+)+)  \s+  ([\d\.]+\,\d\d)`)

func (bhw *BHW) load(content string) error {
	iban, err := bhw.loadIban(content)
	if err != nil {
		return err
	}

	_, found := bhw.accounts[iban]
	if !found {
		bhw.accounts[iban] = accounting.Account{
			IBAN:      iban,
			Name:      fmt.Sprintf("Bausparkonto %s", iban),
			Transfers: []accounting.Transfer{},
		}
	}

	matches := transfersRegex.FindAllStringSubmatch(content, 100)
	if matches == nil {
		return fmt.Errorf("failed to find transfers")
	}
	if len(matches) == 0 {
		return fmt.Errorf("failed to find transfers")
	}

	positiveLength := -1

	for _, match := range matches {
		date := match[1]
		description := match[2]
		amount := match[4]

		if description == "Saldovortrag" {
			positiveLength = len(match[0])
			continue
		}

		if positiveLength == -1 {
			return fmt.Errorf("failed to detect length of positive transaction")
		}
		positive := len(match[0]) == positiveLength

		amount = strings.ReplaceAll(amount, ".", "")
		amount = strings.ReplaceAll(amount, ",", "")
		amountCent, err := strconv.Atoi(amount)
		if err != nil {
			return fmt.Errorf("parse amount: %v", err)
		}

		if !positive {
			amountCent = -amountCent
		}

		t, err := time.Parse("02.01.06", date)
		if err != nil {
			return fmt.Errorf("parse date: %v", err)
		}

		transfer := accounting.Transfer{
			Date:        t,
			Description: description,
			Amount:      amountCent,
		}

		account := bhw.accounts[iban]
		account.Transfers = append(account.Transfers, transfer)
		bhw.accounts[iban] = account
	}

	return nil
}

var ibanRegex = regexp.MustCompile(`IBAN:((\s\S+)+)`)

func (bhw *BHW) loadIban(content string) (string, error) {
	match := ibanRegex.FindStringSubmatch(content)
	if match == nil {
		return "", fmt.Errorf("failed to find IBAN")
	}

	iban := strings.TrimSpace(match[1])

	return iban, nil
}

func (bhw *BHW) Accounts() []accounting.Account {
	accounts := []accounting.Account{}
	for _, account := range bhw.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}
