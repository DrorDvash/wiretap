package config

import (
	"encoding/base64"
	"math/rand"
	"time"
)

const poisonIndex = 2

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Encode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	if len(encoded) > poisonIndex {
		chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		poison := chars[rand.Intn(len(chars))]
		encoded = encoded[:poisonIndex] + string(poison) + encoded[poisonIndex:]
	}
	return encoded
}

func Decode(encoded string) ([]byte, error) {
	if len(encoded) > poisonIndex {
		encoded = encoded[:poisonIndex] + encoded[poisonIndex+1:]
	}
	return base64.StdEncoding.DecodeString(encoded)
}