package cosmos

import (
	"github.com/dan13ram/wpokt-validator/app"
	cosmos "github.com/dan13ram/wpokt-validator/cosmos/client"
)

func ValidateNetwork() {
	_, err := cosmos.NewClient(app.Config.Pocket)
	if err != nil {
		panic(err)
	}
}
