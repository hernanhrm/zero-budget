package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"zb/internal/api"
)

func WriteCurrencies(writer io.Writer, currencies []api.Currency) error {
	if len(currencies) == 0 {
		_, err := fmt.Fprintln(writer, "No currencies found")
		return err
	}

	table := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "CODE\tNAME\tSYMBOL\tDECIMALS"); err != nil {
		return err
	}

	for _, currency := range currencies {
		if _, err := fmt.Fprintf(
			table,
			"%s\t%s\t%s\t%d\n",
			currency.Code,
			currency.Name,
			currency.Symbol,
			currency.DecimalPlaces,
		); err != nil {
			return err
		}
	}

	return table.Flush()
}
