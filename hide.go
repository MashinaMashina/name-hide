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
	maxNameLength = 255
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
		fmt.Println(file.Name())

		if ExtAvailable(filepath.Ext(file.Name())) {
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
		ext := filepath.Ext(file.Name())
		if !ExtAvailable(ext) {
			fmt.Printf("file '%s' has not have extension '%s'\n", file.Name(), LnkExt)
			continue
		}

		filename := strings.TrimSuffix(file.Name(), ext)
		if strings.TrimLeft(filename, string(hideChar)) == "" {
			// файл уже скрыт
			continue
		}

		spaces, err := db.GetUnusedNum()
		if err != nil {
			return fmt.Errorf("generating name: %w", err)
		}

		if spaces+len(ext) > maxNameLength {
			return fmt.Errorf("can not create too long name: %d chars", spaces)
		}

		name := strings.Repeat(string(hideChar), spaces)

		oldPath := fmt.Sprintf("%s/%s", path, file.Name())
		newPath := fmt.Sprintf("%s/%s%s", path, name, ext)

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
