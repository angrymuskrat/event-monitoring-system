package main

import (
	"fmt"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"time"
)

type SelectPostsCountExp struct {
	exps []Experiment
	city string
}

func (c *Connector) runPostHour(hub chan SelectPostsCountExp, times []int64, restarts int) {
	exps, _ := c.postHourRequests(times, restarts)
	hub <- SelectPostsCountExp{exps: exps, city: c.dbName}
	return
}

func (c *Connector) postHourRequests(timestamps []int64, restarts int) (exps []Experiment, err error) {
	for _, timestamp := range timestamps {
		//statement := makeSelectPostsSQL(timestamp, timestamp+60*60)
		statement := makeSelectPostsCountSQL(timestamp, timestamp+60*60)
		start := time.Now()
		exp := Experiment{}
		for i := 0; i < restarts; i++ {
			exp.unitsCount, err = countSelectCount(statement, c.conn)
			if err != nil {
				unilog.Logger().Error("SpeedTestPostHourRequests", zap.Error(err))
				return nil, err
			}
		}
		exp.hour = timestamp
		exp.executionTime = time.Since(start).Microseconds() / int64(restarts)
		exps = append(exps, exp)
	}
	return exps, nil
}

func (exp *SelectPostsCountExp) String() string {
	ans := fmt.Sprintf("dbName: %v:", exp.city)
	eTime, uCounts := int64(0), 0
	for _, exp := range exp.exps {
		eTime += exp.executionTime
		uCounts += exp.unitsCount
	}
	ans += fmt.Sprintf(" time: %v, units: %v, time per unit: %v", eTime, uCounts, float64(eTime)/float64(uCounts))
	return ans
}
