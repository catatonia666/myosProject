package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5431 user=postgres password=world555 dbname=rpg_test sslmode=disable"
	}

	os.Exit(m.Run())
}
