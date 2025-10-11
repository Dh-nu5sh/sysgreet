package bootstrap

import (
	"fmt"
	"io"
)

func logStatus(w io.Writer, action Action, configPath, backupPath string) {
	if w == nil {
		return
	}

	switch action {
	case ActionCreated:
		fmt.Fprintf(w, "sysgreet bootstrap: created default config at %s\n", configPath)
	case ActionKept:
		fmt.Fprintf(w, "sysgreet bootstrap: keeping existing config at %s\n", configPath)
	case ActionOverwritten:
		fmt.Fprintf(w, "sysgreet bootstrap: overwrote config at %s (backup: %s)\n", configPath, backupPath)
	case ActionSkipped:
		fmt.Fprintf(w, "sysgreet bootstrap: skipped bootstrap for %s\n", configPath)
	}
}
