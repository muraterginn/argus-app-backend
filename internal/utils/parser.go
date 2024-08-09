package utils

import (
	"argus-app-backend/internal/models"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

func ParseMessageToJSON(message string) (string, error) {
	msg := models.Message{}

	accountRegex := regexp.MustCompile(`\[ACCOUNT\]\{([^}]*)\}`)
	accountMatch := accountRegex.FindStringSubmatch(message)
	if len(accountMatch) > 1 {
		accountData := parseKeyValue(accountMatch[1])
		msg.Account = models.Account{
			AccountNumber:     int(accountData["account_number"].(float64)),
			AccountBalance:    accountData["account_balance"].(float64),
			AccountDrawdown:   accountData["account_drawdown"].(float64),
			AccountEquity:     accountData["account_equity"].(float64),
			AccountMargin:     accountData["account_margin"].(float64),
			AccountFreeMargin: accountData["account_free_margin"].(float64),
		}
	}

	totalRegex := regexp.MustCompile(`\[TOTAL\]\{([^}]*)\}`)
	totalMatch := totalRegex.FindStringSubmatch(message)
	if len(totalMatch) > 1 {
		totalData := parseKeyValue(totalMatch[1])
		msg.Total = models.Total{
			TotalDD:      int(totalData["total_dd"].(float64)),
			SellDD:       int(totalData["sell_dd"].(float64)),
			BuyDD:        int(totalData["buy_dd"].(float64)),
			NetLot:       totalData["net_lot"].(float64),
			TotalLot:     totalData["total_lot"].(float64),
			TotalSellLot: totalData["total_sell_lot"].(float64),
			TotalBuyLot:  totalData["total_buy_lot"].(float64),
			NetOC:        int(totalData["net_oc"].(float64)),
			TotalOC:      int(totalData["total_oc"].(float64)),
			SellOC:       int(totalData["sell_oc"].(float64)),
			BuyOC:        int(totalData["buy_oc"].(float64)),
			Profit:       totalData["profit"].(float64),
		}
	}

	symbolsRegex := regexp.MustCompile(`\[SYMBOLS\]\{([^}]*)\}`)
	symbolsMatch := symbolsRegex.FindStringSubmatch(message)
	if len(symbolsMatch) > 1 {
		msg.Symbols = parseSymbols(symbolsMatch[1])
	}

	jsonData, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func parseKeyValue(data string) map[string]interface{} {
	result := make(map[string]interface{})
	entries := strings.Split(data, ", ")
	for _, entry := range entries {
		parts := strings.Split(entry, ": ")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = parseValue(value)
		}
	}
	return result
}

func parseSymbols(data string) map[string]models.Symbol {
	result := make(map[string]models.Symbol)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			symbol := strings.TrimSpace(parts[0])
			symbolData := parseKeyValue(parts[1])
			result[symbol] = models.Symbol{
				TotalDD:  int(symbolData["total_dd"].(float64)),
				SellDD:   int(symbolData["sell_dd"].(float64)),
				BuyDD:    int(symbolData["buy_dd"].(float64)),
				NetLot:   symbolData["net_lot"].(float64),
				TotalLot: symbolData["total_lot"].(float64),
				SellLot:  symbolData["sell_lot"].(float64),
				BuyLot:   symbolData["buy_lot"].(float64),
				NetOC:    int(symbolData["net_oc"].(float64)),
				TotalOC:  int(symbolData["total_oc"].(float64)),
				SellOC:   int(symbolData["sell_oc"].(float64)),
				BuyOC:    int(symbolData["buy_oc"].(float64)),
				Profit:   symbolData["profit"].(float64),
			}
		}
	}
	return result
}

func parseValue(value string) interface{} {
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}
	return value
}
