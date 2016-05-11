package main

import (
	"log"
	"strconv"
	"strings"
	"unicode"
	//"testing"
)

// string is "" if this token isn't this
func getDestination(token string) string {
	if token[0] == '@' {
		return token[1:]
	}
	return ""
}

// float is 0 if this token isn't this
func getAmount(token string) float64 {
	log.Println(token)
	res, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0
	}
	return res
}

// float is 0 if this token isn't this
func getPrefixedAmount(token string) (string, float64) {
	var assetCode string
	var amount float64

	runedStr := []rune(token)

	assetCode = symbolToAssetCode(string(runedStr[0]))
	if len(assetCode) == 0 {
		max := len(runedStr) - 1
		assetCode = symbolToAssetCode(string(runedStr[max]))
		if len(assetCode) > 0 {
			amount = getAmount(string(runedStr[:max]))
		}
	} else if len(assetCode) > 0 {
		amount = getAmount(string(runedStr[1:]))
	}

	return assetCode, amount
}

func symbolToAssetCode(symbol string) string {
	switch symbol {
	case "$":
		return "USD"
	case "*":
		return "XLM"
	case "£":
		return "GBP"
	case "€":
		return "EUR"
	case "¥":
		return "JPY"
	}
	return ""
}

func IsAsciiPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func standardPresets(token string) (string, float64) {

	// any 3 or 4 character all upper case string is considered a
	/*
		switch token {
			case "
		}*/

	if len(token) == 3 || len(token) == 4 && IsAsciiPrintable(token) {
		if token == strings.ToUpper(token) {
			return token, 1
		}
	}
	return "", 0
}
