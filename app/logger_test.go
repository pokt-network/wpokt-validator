package app

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestInitLogger(t *testing.T) {
	t.Run("Log level not provided", func(t *testing.T) {
		Config.Logger.Level = ""

		InitLogger()

		assert.Equal(t, log.GetLevel(), log.InfoLevel)
	})

	testCases := []struct {
		level string
		want  log.Level
	}{
		{"debug", log.DebugLevel},
		{"info", log.InfoLevel},
		{"warn", log.WarnLevel},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Log level %s", tc.level), func(t *testing.T) {
			Config.Logger.Level = tc.level

			InitLogger()

			assert.Equal(t, log.GetLevel(), tc.want)
		})
	}

}
