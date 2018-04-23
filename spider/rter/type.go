package rter

import "time"

type CurrencyCode struct {
	Code  string
	Title string
}

type CurrencyRate struct {
	Bank       string
	In         string
	Out        string
	UpdateTime string
}

// CurrencyRateBody json response body
type CurrencyRateBody struct {
	Data [][]string `json:"data"`
}

type News struct {
	Date         time.Time
	Href         string
	ShortContent string
	Title        string
}
