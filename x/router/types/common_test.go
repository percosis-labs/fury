package types_test

import (
	"os"
	"testing"

	"github.com/percosis-labs/fury/app"
)

func TestMain(m *testing.M) {
	app.SetSDKConfig()
	os.Exit(m.Run())
}
