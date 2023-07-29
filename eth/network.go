package eth

import (
	eth "github.com/dan13ram/wpokt-validator/eth/client"
)

func ValidateNetwork() {
	eth.Client.ValidateNetwork()
}
