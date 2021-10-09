package utils

import (
	"math/rand"
	"time"
)

func RandString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[seededRand.Intn(len(letter))]
	}
	return string(b)
}

func RandHexString(n int) string {
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandInt(n int) int {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(n)
}

func RandIntRange(min, max int) int {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(max-min) + min
}

func RandBool() bool {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(2) != 0
}

func RandShuffle(arr []string) []string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	out := append([]string{}, arr...)
	seededRand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}
