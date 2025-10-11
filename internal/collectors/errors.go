package collectors

import (
	"log"
	"os"
)

var debugEnabled = os.Getenv("SYSGREET_DEBUG") != ""

func recordError(context string, err error) {
	if err == nil || !debugEnabled {
		return
	}
	log.Printf("sysgreet %s: %v", context, err)
}
