package main

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const maxLen = 10
const minLen = 5
const maxForInt64 = 10000
const maxLoc = 100

func RandString() string {
	length := seededRand.Int() % (maxLen - minLen) + minLen
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandBool() bool {
	return seededRand.Int() % 2 == 1
}

func RandInt64() int64 {
	return seededRand.Int63() % maxForInt64
}

func RandDouble() float64 {
	return float64(seededRand.Int() % (maxLoc - 1)) + seededRand.Float64();
}

func GeneratePosts(n int) []data.Post {
	posts := *new([]data.Post)
	for i := 0; i < n; i++ {
		posts = append(posts, data.Post{ RandString(), RandString(), RandString(),
			RandBool(), RandString(), RandInt64(), RandInt64(),
			RandInt64(), RandBool(), RandString(), RandString(),
			RandDouble(), RandDouble(), struct{}{}, nil, 0})
	}
	return posts;
}
