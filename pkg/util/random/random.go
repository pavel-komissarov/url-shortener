package random

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode/utf8"
)

var (
	alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_")
)

var ErrInvalidLength = errors.New("stringLength must be > 0")

func NewRandomString(stringLength int) (string, error) {
	if stringLength <= 0 {
		return "", ErrInvalidLength
	}

	var builder strings.Builder

	alphaLength := big.NewInt(int64(utf8.RuneCount(alphabet)))
	for range stringLength {
		n, err := rand.Int(rand.Reader, alphaLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random integer: %w", err)
		}

		builder.WriteByte(alphabet[n.Int64()])
	}

	return builder.String(), nil
}
