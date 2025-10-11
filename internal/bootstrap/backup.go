package bootstrap

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func createBackup(path string, now time.Time) (string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	timestamp := now.UTC().Format("20060102-150405")
	backupName := fmt.Sprintf("%s.bak-%s", base, timestamp)
	backupPath := filepath.Join(dir, backupName)

	if err := os.Rename(path, backupPath); err != nil {
		return "", fmt.Errorf("create backup: %w", err)
	}

	if err := pruneOlderBackups(dir, base, backupName); err != nil {
		return "", err
	}

	return backupPath, nil
}

func pruneOlderBackups(dir, base, keepName string) error {
	pattern := fmt.Sprintf("%s.bak-", base)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("list backups: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, pattern) {
			continue
		}
		if name == keepName {
			continue
		}
		if err := os.Remove(filepath.Join(dir, name)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove old backup %s: %w", name, err)
		}
	}
	return nil
}
