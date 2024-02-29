package accounting

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/golang/snappy"
	"github.com/segmentio/fasthash/fnv1a"
)

type Pdf struct {
	path string
}

func (p *Pdf) Text() (string, error) {
	return p.cached("")
}

func (p *Pdf) Raw() (string, error) {
	return p.cached("-raw")
}

func (p *Pdf) WithLayout() (string, error) {
	return p.cached("-layout")
}

func (p *Pdf) cached(kind string) (string, error) {
	// h, err := md5sum(p.path)
	h, err := fnv1asum(p.path)
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

		b = snappy.Encode(nil, b)

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

	b, err = snappy.Decode(nil, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func fnv1asum(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	h := fnv1a.HashBytes64(b)
	s := fmt.Sprintf("%016x", h)
	return s, nil
}

func md5sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (p *Pdf) pdftotext(kind string) (string, error) {
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
