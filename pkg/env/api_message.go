package env

import (
	"os"
	"strconv"
)

const envPrefix = "TCPTUNNEL_"

var MESSAGE_BATCH_SIZE int

func init() {
	MESSAGE_BATCH_SIZE, err := strconv.Atoi(os.Getenv(envPrefix + "BATCH_SIZE"))
	if err != nil || MESSAGE_BATCH_SIZE == 0 {
		MESSAGE_BATCH_SIZE = 5
	}
}
