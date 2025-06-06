package util

import (
	"strconv"

	"github.com/dan13ram/wpokt-validator/app"
	"github.com/dan13ram/wpokt-validator/models"
)

func UpdateStatusAndConfirmationsForInvalidMint(doc *models.InvalidMint, currentHeight int64) (*models.InvalidMint, error) {
	status := doc.Status
	confirmations, err := strconv.ParseInt(doc.Confirmations, 10, 64)
	if err != nil || confirmations < 0 {
		confirmations = 0
	}

	if status == models.StatusPending || confirmations == 0 {
		status = models.StatusPending
		if app.Config.Pocket.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			mintHeight, err := strconv.ParseInt(doc.Height, 10, 64)
			if err != nil {
				return doc, err
			}
			confirmations = currentHeight - mintHeight
			if confirmations >= app.Config.Pocket.Confirmations {
				status = models.StatusConfirmed
			}
		}
	}

	doc.Status = status
	doc.Confirmations = strconv.FormatInt(confirmations, 10)

	return doc, nil
}

func UpdateStatusAndConfirmationsForBurn(doc *models.Burn, blockNumber int64) (*models.Burn, error) {
	status := doc.Status
	confirmations, err := strconv.ParseInt(doc.Confirmations, 10, 64)
	if err != nil || confirmations < 0 {
		confirmations = 0
	}

	if status == models.StatusPending || confirmations == 0 {
		status = models.StatusPending
		if app.Config.Ethereum.Confirmations == 0 {
			status = models.StatusConfirmed
		} else {
			burnBlockNumber, err := strconv.ParseInt(doc.BlockNumber, 10, 64)
			if err != nil {
				return doc, err
			}

			confirmations = blockNumber - burnBlockNumber
			if confirmations >= app.Config.Ethereum.Confirmations {
				status = models.StatusConfirmed
			}
		}
	}

	doc.Status = status
	doc.Confirmations = strconv.FormatInt(confirmations, 10)
	return doc, nil
}
