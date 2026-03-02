package utils

import (
	"crypto/rand"
	"io"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62From(r io.Reader, n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return "", ErrWrap(err, "failed to generate base 62 from")
	}
	for i, b := range bytes {
		bytes[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(bytes), nil
}

type Base62Generator interface {
	Generate(n int) (string, error)
}

type CryptoBase62Generator struct{}

func (CryptoBase62Generator) Generate(n int) (string, error) {
	return Base62From(rand.Reader, n)
}

type StubGenerator struct {
}

func (s StubGenerator) Generate(n int) (string, error) {
	return alphabet[:n], nil
}
