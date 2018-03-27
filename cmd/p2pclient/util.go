package main

import (
	"math/rand"
	"time"
)

var randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString() string {
	var (
		size = 32
	)
	buffer := make([]byte, size)
	for i := range buffer {
		r := randSrc.Uint32() % uint32(len(letterBytes))
		buffer[i] = letterBytes[r]
	}
	return string(buffer)
}

func decodePayload(b []byte) []byte {
	if len(b) == 0 {
		return b
	}
	return b[1:]
}
