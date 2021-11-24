package main

import (
	"context"
	"flag"
	"fmt"
	utilsrand "github.com/angrymuskrat/event-monitoring-system/utils/rand"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"math/rand"
)

type Connector struct {
	conn   *pgxpool.Pool
	dbName string
	label  string
}

type Experiment struct {
	hour          int64
	executionTime int64
	unitsCount    int
}

func countSelectObjects(statement string, conn *pgxpool.Pool) (count int, err error) {
	count = 0
	rows, err := conn.Query(context.Background(), statement)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		count++
	}
	return count, nil
}

func countSelectCount(statement string, conn *pgxpool.Pool) (count int, err error) {
	row := conn.QueryRow(context.Background(), statement)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewConnector(ctx context.Context, config *Configuration, dbName string, label string) (c *Connector, err error) {
	connConfig, err := pgxpool.ParseConfig(config.makeAuthToken(dbName))
	if err != nil {
		unilog.Logger().Error("unable to parse pg config",
			zap.String("dbName", dbName), zap.String("label", label), zap.Error(err))
		return
	}
	conn, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		unilog.Logger().Error("unable to connect to db",
			zap.String("dbName", dbName), zap.String("label", label), zap.Error(err))
		return
	}
	c = &Connector{
		conn:   conn,
		dbName: dbName,
		label:  label,
	}
	return c, nil
}

func NewListConnectors(ctx context.Context, config *Configuration, dbNames []string, labels []string) (connectors []*Connector, err error) {
	for ind, dbName := range dbNames {
		c, err := NewConnector(ctx, config, dbName, labels[ind])
		if err != nil {
			return nil, err
		}
		connectors = append(connectors, c)
	}
	return connectors, nil
}

func main() {
	configPath := flag.String("sc", "config.toml", "path to configuration file")
	config, err := readConfig(*configPath)
	if err != nil {
		return
	}
	//cities := []string{"moscow", "spb", "nyc", "london", "moscow", "spb", "nyc", "london"}
	dbNames := []string{"moscow", "spb", "nyc", "london"}
	labels := dbNames
	connectors, err := NewListConnectors(context.Background(), &config, dbNames, labels)

	timesCount := 500
	var times []int64 //times := []int64{1572807600, 1582480800}
	seed := int64(12345)
	randomizer := utilsrand.NewFixSeed(seed)
	minTime, maxTime := int64(1514764800), int64(1585699200)

	for i := 0; i < timesCount; i++ {
		times = append(times, randomizer.HourTimestamp(minTime, maxTime))
	}

	cityTimes := make(map[string][]int64)
	rand.NewSource(seed)
	for _, city := range dbNames {
		shuffled := make([]int64, len(times))
		copy(shuffled, times)
		//rand.NewSource(seed * int64(i + 1))
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		//fmt.Println(shuffled[:4], times[:4])
		cityTimes[city] = shuffled
	}

	restarts := 20
	hub := make(chan SelectPostsCountExp)
	for _, conn := range connectors {
		//go conn.runSpeedTestPostHour(hub, times, restarts)
		go conn.runPostHour(hub, cityTimes[conn.dbName], restarts)
	}

	fmt.Println("Parallel test")
	testsCount := 0
	for test := range hub {
		fmt.Println(test)
		testsCount += 1
		if testsCount == len(connectors) {
			break
		}
	}
	fmt.Println("Sequence test")

	for _, conn := range connectors {
		//go conn.runSpeedTestPostHour(hub, times, restarts)
		go conn.runPostHour(hub, cityTimes[conn.dbName], restarts)
		test := <-hub
		fmt.Println(test)
	}

}
