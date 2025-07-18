package display

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"math/big"
	"os"
)

type TablePrinter struct{}

func NewTablePrinter() *TablePrinter {
	return &TablePrinter{}
}

func (p *TablePrinter) Print(address string, assets []AssetInfo, totalUSD *big.Float) {
	fmt.Printf("Address: %s\n\n", address)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"TOKEN", "BALANCE", "VALUE (USD)"})

	for _, a := range assets {
		table.Append([]string{
			a.Symbol,
			a.Balance.Text('f', 6),
			fmt.Sprintf("$%.2f", a.USDValue),
		})
	}

	table.Append([]string{"TOTAL", "", fmt.Sprintf("$%.2f", totalUSD)})
	table.Render()
}
