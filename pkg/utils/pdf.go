package utils

import (
	"bytes"
	"github.com/ledongthuc/pdf"
)

func GetPdfText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(b); err != nil {
		return "", err
	}
	return buf.String(), nil
}
