package util

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ParseMessageSenderEvent(
	events []abci.Event,
) (string, error) {
	for _, event := range events {
		if strings.EqualFold(event.Type, "message") {
			for _, attr := range event.Attributes {
				if strings.EqualFold(string(attr.Key), "sender") {
					sender := string(attr.Value)
					return sender, nil
				}
			}
		}
	}
	return "", fmt.Errorf("no sender found in message events")
}

func ParseCoinsReceivedEvents(
	denom string,
	receiver string,
	events []abci.Event,
) (sdk.Coin, error) {
	total := sdk.NewCoin(denom, math.NewInt(0))
	for _, event := range events {
		if strings.EqualFold(event.Type, "coin_received") {
			for _, attr := range event.Attributes {
				if strings.EqualFold(string(attr.Key), "receiver") && strings.EqualFold(string(attr.Value), receiver) {
					for _, attr := range event.Attributes {
						if strings.EqualFold(string(attr.Key), "amount") {
							amountStr := string(attr.Value)
							amount, err := sdk.ParseCoinNormalized(amountStr)
							if err != nil {
								return total, fmt.Errorf("unable to parse coin amount: %v", err)
							}
							if amount.Denom != denom {
								return total, fmt.Errorf("invalid coin denom: %s", amount.Denom)
							}
							total = total.Add(amount)
						}
					}
				}
			}
		}
	}
	return total, nil
}

func ParseCoinsSpentEvents(
	denom string,
	events []abci.Event,
) (string, sdk.Coin, error) {
	total := sdk.NewCoin(denom, math.NewInt(0))
	spender := ""
	for _, event := range events {
		if strings.EqualFold(event.Type, "coin_spent") {
			for _, attr := range event.Attributes {
				if strings.EqualFold(string(attr.Key), "spender") {
					newSpender := string(attr.Value)
					if spender != "" && !strings.EqualFold(spender, newSpender) {
						return spender, total, fmt.Errorf("multiple spenders found in coin spent events")
					}
					spender = newSpender
				}
				if strings.EqualFold(string(attr.Key), "amount") {
					amountStr := string(attr.Value)
					amount, err := sdk.ParseCoinNormalized(amountStr)
					if err != nil {
						return spender, total, err
					}
					if amount.Denom != denom {
						return spender, total, fmt.Errorf("invalid coin denom: %s", amount.Denom)
					}
					total = total.Add(amount)
				}
			}
		}
	}
	return spender, total, nil
}
