package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
)

type Response struct {
	Rates map[string]float64 `json:"rates"`
}

func main() {

	ammount, fromCurrency, toCurrency, newAmmount := createForm()

	rate := fetchData(fromCurrency, toCurrency)

	finalConversion := newAmmount * float64(rate)

	fmt.Println("O valor", ammount, fromCurrency, "fica", finalConversion, toCurrency)

}

func createForm() (string, string, string, float64) {

	var ammount string
	var fromCurrency string
	var toCurrency string

	form := huh.NewForm(
		huh.NewGroup(

			huh.NewInput().
				Title("Ammount").
				Value(&ammount),

			huh.NewSelect[string]().
				Title("From Currency").
				Options(
					huh.NewOption("EUR", "EUR"),
					huh.NewOption("USD", "USD"),
					huh.NewOption("GBP", "GBP"),
				).
				Value(&fromCurrency),

			huh.NewSelect[string]().
				Title("To Currency").
				Options(
					huh.NewOption("EUR", "EUR"),
					huh.NewOption("USD", "USD"),
					huh.NewOption("GBP", "GBP"),
				).
				Value(&toCurrency),
		),
	)

	form.Run()

	//Verifica se estamos a tentar converter para a mesma currency
	if fromCurrency == toCurrency {
		fmt.Println("Currency nao pode ser a mesma")
		os.Exit(1)
	}

	//Transforma o valor da qauntidade introduzida para float 64
	newAmmount, err := strconv.ParseFloat(ammount, 64)
	if err != nil {
		panic(err)
	}

	return ammount, fromCurrency, toCurrency, newAmmount

}

func fetchData(fromCurrency string, toCurrency string) float64 {
	var data Response

	url := fmt.Sprintf(
		"https://api.frankfurter.app/latest?from=%s&to=%s",
		fromCurrency,
		toCurrency,
	)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(body))

	json.Unmarshal(body, &data)

	rate := data.Rates[toCurrency]

	return rate
}
