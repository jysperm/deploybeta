package utils

import (
	"crypto/rand"
	"encoding/binary"
	mathRand "math/rand"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *mathRand.Rand

func RandomString(length int) string {
	buffer := make([]byte, length)

	for i := range buffer {
		buffer[i] = base62Charset[seededRand.Intn(len(base62Charset))]
	}

	return string(buffer)
}

func init() {
	buffer := make([]byte, 8)

	_, err := rand.Read(buffer)

	if err != nil {
		panic(err)
	}

	seededRand = mathRand.New(mathRand.NewSource(int64(binary.LittleEndian.Uint64(buffer))))
}
