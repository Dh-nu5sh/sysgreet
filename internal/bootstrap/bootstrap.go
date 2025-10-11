package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"
)

// ErrUserCancelled indicates the user chose to cancel during an interactive prompt.
var ErrUserCancelled = errors.New("sysgreet bootstrap cancelled by user")

// Action describes what happened during bootstrap.
type Action string

const (
	ActionSkipped     Action = "skipped"
	ActionCreated     Action = "created"
	ActionOverwritten Action = "overwritten"
	ActionKept        Action = "kept"
)

// IO bundles input/output writers used during bootstrap operations.
type IO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Options struct {
	FlagPolicy  string
	EnvPolicy   string
	Interactive bool
}

// Result captures the outcome of the bootstrap flow.
type Result struct {
	Action     Action
	ConfigPath string
	BackupPath string
	Policy     PolicyValue
	Prompted   bool
}

func normalizeIO(ioCfg IO) IO {
	if ioCfg.Stdin == nil {
		ioCfg.Stdin = strings.NewReader("")
	}
	if ioCfg.Stdout == nil {
		ioCfg.Stdout = io.Discard
	}
	if ioCfg.Stderr == nil {
		ioCfg.Stderr = io.Discard
	}
	return ioCfg
}

// Bootstrap ensures the sysgreet configuration exists according to policy.
func Bootstrap(ctx context.Context, cfgPath string, ioCfg IO, opts Options) (Result, error) {
	ioCfg = normalizeIO(ioCfg)

	if cfgPath == "" {
		return Result{}, fmt.Errorf("bootstrap: config path required")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	resolution, err := ResolvePolicy(opts.FlagPolicy, opts.EnvPolicy, opts.Interactive)
	if err != nil {
		return Result{}, err
	}

	now := time.Now()

	result := Result{
		ConfigPath: cfgPath,
		Policy:     resolution.Value,
	}

	info, statErr := os.Stat(cfgPath)
	if statErr == nil {
		if info.IsDir() {
			return result, fmt.Errorf("bootstrap: config path %s is a directory", cfgPath)
		}

		switch resolution.Value {
		case PolicyKeep:
			result.Action = ActionKept
			logStatus(ioCfg.Stderr, result.Action, cfgPath, "")
			return result, nil
		case PolicyOverwrite:
			if err := ctx.Err(); err != nil {
				return result, err
			}
			data, err := renderDefaultConfig(now)
			if err != nil {
				return result, fmt.Errorf("bootstrap: render default config: %w", err)
			}
			backupPath, err := createBackup(cfgPath, now)
			if err != nil {
				return result, fmt.Errorf("bootstrap: backup existing config: %w", err)
			}
			if err := AtomicWriteFile(cfgPath, data, 0o644); err != nil {
				return result, fmt.Errorf("bootstrap: write default config: %w", err)
			}
			result.Action = ActionOverwritten
			result.BackupPath = backupPath
			logStatus(ioCfg.Stderr, result.Action, cfgPath, backupPath)
			return result, nil
		case PolicyPrompt:
			if err := ctx.Err(); err != nil {
				return result, err
			}
			outcome, err := PromptForOverwrite(ioCfg, cfgPath)
			if err != nil {
				return result, err
			}
			result.Prompted = true
			switch outcome.Decision {
			case PromptKeep:
				result.Action = ActionKept
				logStatus(ioCfg.Stderr, result.Action, cfgPath, "")
				return result, nil
			case PromptOverwrite:
				data, err := renderDefaultConfig(now)
				if err != nil {
					return result, fmt.Errorf("bootstrap: render default config: %w", err)
				}
				backupPath, err := createBackup(cfgPath, now)
				if err != nil {
					return result, fmt.Errorf("bootstrap: backup existing config: %w", err)
				}
				if err := AtomicWriteFile(cfgPath, data, 0o644); err != nil {
					return result, fmt.Errorf("bootstrap: write default config: %w", err)
				}
				result.Action = ActionOverwritten
				result.BackupPath = backupPath
				logStatus(ioCfg.Stderr, result.Action, cfgPath, backupPath)
				return result, nil
			case PromptCancel:
				result.Action = ActionSkipped
				return result, ErrUserCancelled
			default:
				return result, fmt.Errorf("bootstrap: unknown prompt decision %s", outcome.Decision)
			}
		default:
			return result, fmt.Errorf("bootstrap: unsupported policy %s", resolution.Value)
		}
	}
	if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
		return result, fmt.Errorf("bootstrap: stat config: %w", statErr)
	}

	if err := ctx.Err(); err != nil {
		return result, err
	}

	if resolution.Value == PolicyKeep {
		result.Action = ActionKept
		logStatus(ioCfg.Stderr, result.Action, cfgPath, "")
		return result, nil
	}

	data, err := renderDefaultConfig(now)
	if err != nil {
		return result, fmt.Errorf("bootstrap: render default config: %w", err)
	}

	if err := AtomicWriteFile(cfgPath, data, 0o644); err != nil {
		return result, fmt.Errorf("bootstrap: write default config: %w", err)
	}

	result.Action = ActionCreated
	logStatus(ioCfg.Stderr, result.Action, cfgPath, "")
	return result, nil
}
