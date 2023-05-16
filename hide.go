package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	hideChar      = rune(160)
	maxNameLength = 251
	LnkExt        = ".lnk"
)

func hide(c *cli.Context) error {
	path, err := filepath.Abs(c.Path("path"))
	if err != nil {
		return fmt.Errorf("invalid path '%s': %w", c.Path("path"), err)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("list directory: %w", err)
	}

	count := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == LnkExt {
			count++
		}
	}

	if count == 0 {
		return fmt.Errorf("not found links in folder '%s'", path)
	}

	// не сможем всем раздать новые имена
	if count > maxNameLength {
		return fmt.Errorf("too many links")
	}

	db := NewDatabase(path)
	if err = db.Init(); err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	for _, file := range files {
		if filepath.Ext(file.Name()) != LnkExt {
			continue
		}

		// файл уже скрыт
		if strings.TrimLeft(file.Name(), string(hideChar)) == LnkExt {
			continue
		}

		spaces, err := db.GetUnusedNum()
		if err != nil {
			return fmt.Errorf("generating name: %w", err)
		}

		if spaces > maxNameLength {
			return fmt.Errorf("can not create too long name: %d chars", spaces)
		}

		name := strings.Repeat(string(hideChar), spaces)

		oldPath := fmt.Sprintf("%s/%s", path, file.Name())
		newPath := fmt.Sprintf("%s/%s.lnk", path, name)

		if err = os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("rename '%s': %w", file.Name(), err)
		}

		if err = db.SaveName(file.Name(), spaces); err != nil {
			return fmt.Errorf("save name '%s' as %d: %w", file.Name(), spaces, err)
		}

		fmt.Printf("%s renamed to %d spaces\n", file.Name(), spaces)
	}

	return nil
}
