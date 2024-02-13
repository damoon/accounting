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
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	txts, err := accounting.ListTxts("data")
	if err != nil {
		return fmt.Errorf("listing txt files: %v", err)
	}

	bhw := bhw.NewBhw()

	for _, txt := range txts {
		b, err := os.ReadFile(txt)
		if err != nil {
			return err
		}

		content := string(b)

		err = bhw.LoadFrom(content)
		if err != nil {
			return err
		}
	}

	/*
		tbl := table.New("ISIN", "Name")
		for _, account := range bhw.Accounts() {
			tbl.AddRow(account.IBAN, account.Name)
		}
		tbl.Print()

		tbl = table.New("Date", "Description", "Amount") //, "Source", "Target")
		for _, transfer := range bhw.Transfers() {
			tbl.AddRow(transfer.Date, transfer.Description, transfer.Amount) //, transfer.Source, transfer.Target)
		}
		tbl.Print()
	*/

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SortBy([]table.SortBy{
		{Name: "ISIN", Mode: table.Asc},
	})
	t.AppendHeader(table.Row{"ISIN", "Name"})
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
	t.AppendHeader(table.Row{"Date", "DateRaw", "Description", "Amount"})
	for _, transfer := range bhw.Transfers() {
		t.AppendRow(table.Row{transfer.Date, transfer.Date, transfer.Description, transfer.Amount})
	}
	t.Render()

	return nil
}
