package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Convert(input string) (string, error) {
	text := strings.TrimSpace(input)
	if text == "" {
		return "", errors.New("пустая строка")
	}

	for _, ch := range text {
		if strings.ContainsRune(".- /", ch) == false {
			return morse.ToMorse(text), nil
		}
	}

	return morse.ToText(text), nil
}
