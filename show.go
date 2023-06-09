package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/urfave/cli/v2"
)

func show(c *cli.Context) error {
	path, err := filepath.Abs(c.Path("path"))
	if err != nil {
		return fmt.Errorf("invalid path '%s': %w", c.Path("path"), err)
	}

	db := NewDatabase(path)
	if !db.Exists() {
		return fmt.Errorf("database file do not exists in '%s'", db.FileName())
	}
	if err = db.Init(); err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	list, err := db.List()
	if err != nil {
		return fmt.Errorf("get names list: %w", err)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("list directory: %w", err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if !ExtAvailable(ext) {
			continue
		}

		filename := strings.TrimSuffix(file.Name(), ext)
		if strings.TrimLeft(filename, string(hideChar)) != "" {
			// В имени файла посторонние символы - файл не скрывался
			continue
		}

		spaces := utf8.RuneCountInString(file.Name()) - len(ext)

		name, err := list.GetName(spaces)
		if err != nil {
			fmt.Printf("get original name by %d spaces: %w\n", spaces, err)
		}

		oldPath := fmt.Sprintf("%s/%s", path, file.Name())
		newPath := fmt.Sprintf("%s/%s", path, name)

		if err = os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("rename '%s': %w", file.Name(), err)
		}

		if err = db.FreeNum(spaces); err != nil {
			return fmt.Errorf("mark spaces as free: %w", err)
		}

		fmt.Printf("%d spaces renamed to %s\n", spaces, name)
	}

	return db.FreeSpace()
}
