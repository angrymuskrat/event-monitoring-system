package main

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const maxForInt64 = 10000
const maxLoc = 100

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandBool() bool {
	return seededRand.Int() % 2 == 1
}

func RandInt64(maxV int64) int64 {
	return seededRand.Int63() % maxV
}

func RandDouble() float64 {
	return float64(seededRand.Int() % (maxLoc - 1)) + seededRand.Float64();
}

func GeneratePosts(n int) []data.Post {
	posts := *new([]data.Post)
	for i := 0; i < n; i++ {
		posts = append(posts, data.Post{ ID: RandString(20), Shortcode: RandString(10), ImageURL: RandString(30),
			IsVideo: RandBool(), Caption: RandString(100), CommentsCount: RandInt64(1000), Timestamp: RandInt64(1000),
			LikesCount: RandInt64(1000), IsAd: RandBool(), AuthorID: RandString(15), LocationID: RandString(15),
			Lat: RandDouble(), Lon: RandDouble()})
	}
	return posts;
}

