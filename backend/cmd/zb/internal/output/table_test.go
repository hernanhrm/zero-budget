package output

import (
	"bytes"
	"strings"
	"testing"

	"zb/internal/api"
)

func TestWriteCurrencies(t *testing.T) {
	tests := []struct {
		name       string
		currencies []api.Currency
		want       []string
	}{
		{
			name: "renders currency rows",
			currencies: []api.Currency{
				{
					Code:          "USD",
					Name:          "US Dollar",
					Symbol:        "$",
					DecimalPlaces: 2,
				},
			},
			want: []string{"CODE", "USD", "US Dollar"},
		},
		{
			name:       "renders empty state",
			currencies: nil,
			want:       []string{"No currencies found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := WriteCurrencies(&buffer, tt.currencies); err != nil {
				t.Fatalf("WriteCurrencies() error = %v", err)
			}

			output := buffer.String()
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Fatalf("output %q does not contain %q", output, want)
				}
			}
		})
	}
}
