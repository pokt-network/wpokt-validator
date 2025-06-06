package util

import (
	"fmt"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TransferEvent struct {
	Sender   string
	Receiver string
	Amount   sdk.Coin
}

func FindAttributeValue(
	events []abci.EventAttribute,
	key string,
) (string, error) {
	for _, event := range events {
		if strings.EqualFold(string(event.Key), key) {
			return string(event.Value), nil
		}
	}
	return "", fmt.Errorf("no attribute found with key: %s", key)
}

func ParseTransferEvents(
	events []abci.Event,
	recipient string,
	denom string,
) ([]TransferEvent, error) {
	transfers := []TransferEvent{}
	for _, event := range events {
		if strings.EqualFold(event.Type, "transfer") {
			sender, err := FindAttributeValue(event.Attributes, "sender")
			if err != nil {
				return transfers, err
			}
			receiver, err := FindAttributeValue(event.Attributes, "recipient")
			if err != nil {
				return transfers, err
			}
			if !strings.EqualFold(receiver, recipient) {
				// skip any transfers that are not to the recipient
				continue
			}
			amountStr, err := FindAttributeValue(event.Attributes, "amount")
			if err != nil {
				return transfers, err
			}
			amount, err := sdk.ParseCoinNormalized(amountStr)
			if err != nil {
				return transfers, err
			}
			if amount.Denom != denom {
				// skip any transfers that are not in the correct denom
				continue
			}
			transfers = append(transfers, TransferEvent{
				Sender:   sender,
				Receiver: receiver,
				Amount:   amount,
			})
		}
	}
	return transfers, nil
}
