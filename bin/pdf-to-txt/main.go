package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	force := os.Getenv("LIST_ALL_PDFS") != ""

	err := run(force)
	if err != nil {
		log.Fatal(err)
	}
}

func run(force bool) error {
	files, err := getPdfs("data")
	if err != nil {
		return err
	}

	for _, f := range files {
		base := filepath.Base(f)
		base, _ = strings.CutSuffix(base, ".pdf")
		base = fmt.Sprintf("%s.txt", base)
		dir := filepath.Dir(f)
		txtPath := path.Join(dir, base)

		if _, err := os.Stat(txtPath); err == nil {
			if !force {
				continue
			}
		}

		fmt.Printf("%v\n", f)
	}
	return nil
}

func getPdfs(rootpath string) ([]string, error) {
	list := []string{}

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".pdf" {
			list = append(list, path)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return list, nil
}
