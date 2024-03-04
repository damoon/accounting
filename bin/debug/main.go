package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/damoon/accounting"
	"github.com/damoon/accounting/bhw"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/pkg/profile"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	switch os.Getenv("PROFILE") {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	case "trace":
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	}

	pdfs, err := accounting.Pdfs("data")
	if err != nil {
		return fmt.Errorf("listing files: %v", err)
	}

	bhw := bhw.NewBhw()

	for _, pdf := range pdfs {
		log.Printf("processing %v", pdf)

		err = bhw.LoadFrom(pdf)
		if err != nil {
			return err
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SortBy([]table.SortBy{
		{Name: "IBAN", Mode: table.Asc},
	})
	t.AppendHeader(table.Row{"IBAN", "Name"})
	for _, account := range bhw.Accounts() {
		t.AppendRow(table.Row{account.IBAN, account.Name})
	}
	t.Render()

	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SortBy([]table.SortBy{
		{Name: "DateRaw", Mode: table.Asc},
		{Name: "Description", Mode: table.Asc},
	})

	dateTransformer := text.Transformer(func(val interface{}) string {
		return val.(time.Time).Format("02.01.2006")
	})
	amountTransformer := text.Transformer(func(val interface{}) string {
		v := val.(int)
		return fmt.Sprintf("%.2f", float64(v)/100)
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "Date",
			Transformer: dateTransformer,
		},
		{
			Name:   "DateRaw",
			Hidden: true,
		},
		{
			Name:        "Amount",
			Align:       text.AlignRight,
			Transformer: amountTransformer,
		},
	})
	t.AppendHeader(table.Row{"IBAN", "Date", "DateRaw", "Description", "Amount"})

	accounts := bhw.Accounts()
	for _, account := range accounts {
		for _, transfer := range account.Transfers {
			t.AppendRow(table.Row{account.IBAN, transfer.Date, transfer.Date, transfer.Description, transfer.Amount})
		}
	}
	t.Render()

	return nil
}
