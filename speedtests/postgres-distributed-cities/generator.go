package main

import (
	utilsRand "github.com/angrymuskrat/event-monitoring-system/utils/rand"
	//posRand "github.com/angrymuskrat/event-monitoring-system/utils/rand/positional"
	"math/rand"
)

type Generator struct {
	seed       int64
	simpleRand *utilsRand.SimpleRand
	mathRand   *rand.Rand
	//positionalGen *posRand.Rand
}

func NewGenerator(seed int64) (gen *Generator) {
	sRand := utilsRand.NewFixSeed(seed)
	mRand := rand.New(rand.NewSource(seed))
	return &Generator{simpleRand: sRand, mathRand: mRand}
}

func (g *Generator) genHours(timesCount int, minTime, maxTime int64) (times []int64) {
	for i := 0; i < timesCount; i++ {
		times = append(times, g.simpleRand.HourTimestamp(minTime, maxTime))
	}
	return times
}
