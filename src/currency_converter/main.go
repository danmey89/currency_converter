package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ExRate struct {
	Rates map[string]float64 `json:"conversion_rates"`
}

func getKey() string {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading.env file: %s", err)
	}
	key := os.Getenv("EXCHANGERATE_API")
	return key
}

func getConversion(baseCurrency string, targetCurrency string, input float64) (float64, float64) {

	key := getKey()
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", key, baseCurrency)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response: %s", err)
	}
	var exchangeRate ExRate

	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		log.Fatalf("Error unmarshalling: %s", err)
	}
	rate := exchangeRate.Rates[targetCurrency]
	
	output := input * rate
	return rate, output
}

func main() {
	
	baseCurrency := "EUR"
	targetCurrency := "CAD"
	input, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		log.Fatalf("Input error: amount value not numeric")
	}
	
	rate, out := getConversion(baseCurrency,targetCurrency,input)
	
	fmt.Printf("Converted value %f %s\nConversion rate: %f\n", out, targetCurrency, rate)
}
