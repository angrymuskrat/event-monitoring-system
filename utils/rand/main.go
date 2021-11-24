package rand

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type SimpleRand struct {
	seeded *rand.Rand
}

func New() *SimpleRand {
	r := SimpleRand{}
	r.seeded = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &r
}

func NewFixSeed(seed int64) *SimpleRand {
	r := SimpleRand{}
	r.seeded = rand.New(rand.NewSource(seed))
	return &r
}

func (r *SimpleRand) Point(point data.Point, delta data.Point) data.Point {
	return data.Point{Lat: point.Lat + r.DeltaDouble(delta.Lat), Lon: point.Lon + r.DeltaDouble(delta.Lon)}
}

func (r *SimpleRand) FixString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.seeded.Intn(len(charset))]
	}
	return string(b)
}

func (r *SimpleRand) AbsInt64(min, max int64) int64 {
	return int64(r.seeded.Uint64()>>1)%(max-min) + min
}

func (r *SimpleRand) HourTimestamp(min, max int64) int64 {
	timestamp := r.AbsInt64(min, max)
	timestamp = timestamp / 3600 * 3600
	return timestamp
}

func (r *SimpleRand) String(min, max int) string {
	length := int(r.AbsInt64(int64(min), int64(max)))
	return r.FixString(length)
}

func (r *SimpleRand) Bool() bool {
	return r.seeded.Int()%2 == 1
}

func (r *SimpleRand) Sign() float64 {
	return float64(1 + -2*r.seeded.Int()%2)
}

func (r *SimpleRand) Double() float64 {
	return r.seeded.Float64()
}

func (r *SimpleRand) DeltaDouble(delta float64) float64 {
	return r.Sign() * r.Double() * delta
}
