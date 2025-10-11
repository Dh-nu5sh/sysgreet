package collectors

import (
	"log"
	"os"
)

var debugEnabled = os.Getenv("HOSTINFO_DEBUG") != ""

func recordError(context string, err error) {
	if err == nil || !debugEnabled {
		return
	}
	log.Printf("hostinfo %s: %v", context, err)
}
