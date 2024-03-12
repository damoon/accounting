package bhw

import (
	_ "embed"
	"reflect"
	"testing"
	"time"

	"github.com/damoon/accounting"
)

//go:embed tests/file1.txt
var file1 string

//go:embed tests/file2.txt
var file2 string

type TestFile struct {
	content string
}

func (t TestFile) WithLayout() (string, error) {
	return t.content, nil
}

func mustParse(date string) time.Time {
	t, err := time.Parse("02.01.06", date)
	if err != nil {
		panic(err)
	}

	return t
}

func TestBHW_LoadFrom(t *testing.T) {
	tests := []struct {
		name    string
		pdfs    []TestFile
		wantErr bool
		want    []accounting.Account
	}{
		{
			name: "Empty",
			pdfs: []TestFile{},
			want: []accounting.Account{},
		},
		{
			name: "2 Pdfs",
			pdfs: []TestFile{
				{file1},
				{file2},
			},
			want: []accounting.Account{
				{
					IBAN: "DE12 3456 7890 1234 5678 90",
					Name: "Bausparkonto DE12 3456 7890 1234 5678 90",
					Transfers: []accounting.Transfer{
						{Date: mustParse("31.01.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("28.02.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("29.03.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.04.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.05.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("28.06.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.07.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.08.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.09.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.10.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("29.11.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.12.19"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.12.19"), Description: "Zinsen auf Bausparguthaben", Amount: 912},
						{Date: mustParse("31.12.19"), Description: "Abgeltungsteuer", Amount: -228},
						{Date: mustParse("31.12.19"), Description: "Solidaritätszuschlag", Amount: -12},

						{Date: mustParse("31.01.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("28.02.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.03.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("28.04.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.05.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.06.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.07.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.08.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("29.09.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("31.10.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.11.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("29.12.23"), Description: "Einzahlung durch Lastschrifteinzug", Amount: 30000},
						{Date: mustParse("30.12.23"), Description: "Guthabenzinsen", Amount: 2358},
						{Date: mustParse("30.12.23"), Description: "Abgeltungsteuer", Amount: -590},
						{Date: mustParse("30.12.23"), Description: "Solidaritätszuschlag", Amount: -32},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bhw := NewBhw()
			for _, pdf := range tt.pdfs {
				if err := bhw.LoadFrom(pdf); (err != nil) != tt.wantErr {
					t.Errorf("BHW.LoadFrom() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if got := bhw.Accounts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BHW.Accounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
