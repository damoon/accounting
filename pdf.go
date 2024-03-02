package accounting

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Pdf struct {
	path string
}

func (p Pdf) String() string {
	return p.path
}

func (p Pdf) Text() (string, error) {
	return p.cached("")
}

func (p Pdf) Raw() (string, error) {
	return p.cached("-raw")
}

func (p Pdf) WithLayout() (string, error) {
	return p.cached("-layout")
}

func (p Pdf) cached(kind string) (string, error) {
	h, err := p.hash()
	if err != nil {
		return "", err
	}

	kind_ := kind
	if kind_ == "" {
		kind_ = "default"
	}

	cacheDir := filepath.Join(".cache", "pdftotext", h[0:2])

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s-%s", h, kind_))

	if _, err := os.Stat(cachePath); err == nil {
		// exists
	} else if errors.Is(err, os.ErrNotExist) {
		// not cached
		txt, err := p.pdftotext(kind)
		if err != nil {
			return "", err
		}

		b := []byte(txt)

		err = os.WriteFile(cachePath, b, os.ModePerm)
		if err != nil {
			return "", err
		}

		return txt, nil
	} else {
		// some error
		return "", err
	}

	b, err := os.ReadFile(cachePath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (p Pdf) hash() (string, error) {
	hash := md5.New()

	_, err := hash.Write([]byte(p.path))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (p Pdf) pdftotext(kind string) (string, error) {
	cmd := exec.Command("pdftotext", kind, p.path, "/dev/stdout")
	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func Pdfs(rootpath string) ([]Pdf, error) {
	list := []Pdf{}

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".pdf" {
			pdf := Pdf{
				path: path,
			}
			list = append(list, pdf)
		}

		return nil
	})

	if err != nil {
		return []Pdf{}, err
	}

	return list, nil
}
