package pokt

import (
	pokt "github.com/dan13ram/wpokt-validator/pokt/client"
)

func ValidateNetwork() {
	pokt.Client.ValidateNetwork()
}
