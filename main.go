package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ExRate struct {
	Rates map[string]float64 `json:"conversion_rates"`
}

func getKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading.env file: %s", err)
	}
	key := os.Getenv("EXCHANGERATE_API")
	return key
}

func getConversion(baseCurrency string, targetCurrency string, input float64) (float64, float64) {

	var exchangeRate ExRate
	key := getKey()
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", key, baseCurrency)

	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response: %s", err)
	}
	err = json.Unmarshal(body, &exchangeRate)
	if err != nil {
		log.Fatalf("Error unmarshalling: %s", err)
	}
	output := input * exchangeRate.Rates[targetCurrency]
	return exchangeRate.Rates[targetCurrency], output
}

func validateIn(args []string) error {
	if len(args) == 4 {
		return nil
	}
	if len(args) == 1 {
		fmt.Println("Please enter amount, original currency and target currency")
		os.Exit(1)
	}
	return errors.New("Invalid number of Arguments, please enter exactly three arguments.\nExpected arguments: <amount> <original currency> <target currency>")
}

func validateAmt(amt string) (float64, error) {
	input, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		return 0, errors.New("Invalid amount, please enter a valid numeric value")
	}
	return input, nil
}

func validateCur(from string, to string) (string, string, error) {
	reg := regexp.MustCompile("^[a-zA-z]{3}$")
	if reg.Match([]byte(from)) || reg.Match([]byte(to)) {
		return strings.ToUpper(from), strings.ToUpper(to), nil
	}
	return "", "", errors.New("Invalid currency format: please enter valid three letter currency code")
}

func main() {
	err := validateIn(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	baseCurrency, targetCurrency, err := validateCur(os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	amount, err := validateAmt(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	rate, out := getConversion(baseCurrency, targetCurrency, amount)

	fmt.Printf("Converted value %f %s\nConversion rate: %f\n", out, targetCurrency, rate)
}
