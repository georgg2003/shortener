package utils

import (
	"crypto/rand"
	"io"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62From(r io.Reader, n int) string {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(r, bytes); err != nil {
		panic(err)
	}
	for i, b := range bytes {
		bytes[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(bytes)
}

type Base62Generator interface {
	Generate(n int) string
}

type CryptoBase62Generator struct{}

func (CryptoBase62Generator) Generate(n int) string {
	return Base62From(rand.Reader, n)
}

type StubGenerator struct {
}

func (s StubGenerator) Generate(n int) string {
	return alphabet[:n]
}
