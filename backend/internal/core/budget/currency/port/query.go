package port

type Currency struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int16  `json:"decimalPlaces"`
}
