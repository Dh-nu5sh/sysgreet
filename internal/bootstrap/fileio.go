package bootstrap

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// AtomicWriteFile writes data to the given path using a same-directory temp file and rename.
func AtomicWriteFile(path string, data []byte, perm fs.FileMode) error {
	if path == "" {
		return fmt.Errorf("atomic write: empty path")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("atomic write: mkdir %s: %w", dir, err)
	}

	if perm == 0 {
		perm = 0o644
	}

	tmp, err := os.CreateTemp(dir, ".sysgreet-*.tmp")
	if err != nil {
		return fmt.Errorf("atomic write: create temp: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}

	if _, err := tmp.Write(data); err != nil {
		cleanup()
		return fmt.Errorf("atomic write: write temp: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		cleanup()
		return fmt.Errorf("atomic write: sync temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("atomic write: close temp: %w", err)
	}
	if err := os.Chmod(tmpName, perm); err != nil {
		cleanup()
		return fmt.Errorf("atomic write: chmod temp: %w", err)
	}

	if err := replaceFile(tmpName, path); err != nil {
		cleanup()
		return fmt.Errorf("atomic write: replace: %w", err)
	}

	return nil
}

func replaceFile(tmpPath, finalPath string) error {
	if err := os.Rename(tmpPath, finalPath); err == nil {
		return nil
	}

	// Windows cannot replace an existing file with os.Rename; remove and retry.
	if err := os.Remove(finalPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.Rename(tmpPath, finalPath)
}
